package main

import (
	"bufio"
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"dotfiles/src/helpers"
	"dotfiles/src/utils"

	"github.com/logrusorgru/aurora/v4"
)

type snapshot struct {
	ID       string    `json:"id"`
	Time     time.Time `json:"time"`
	Hostname string    `json:"hostname"`
}

type entry struct {
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	RequestID string `json:"requestId"`
	Message   struct {
		ID    string `json:"id"`
		Model string `json:"model"`
		Usage *struct {
			InputTokens              int64 `json:"input_tokens"`
			OutputTokens             int64 `json:"output_tokens"`
			CacheCreationInputTokens int64 `json:"cache_creation_input_tokens"`
			CacheReadInputTokens     int64 `json:"cache_read_input_tokens"`
		} `json:"usage"`
	} `json:"message"`
}

type tokens struct {
	Requests   int64
	Input      int64
	Output     int64
	CacheWrite int64
	CacheRead  int64
}

type response struct {
	Device string
	Model  string
	Month  string
	Tokens tokens
}

func main() {
	if !utils.IsCommandInPath("restic") {
		fmt.Println(aurora.Red("Error: restic not installed"))
		os.Exit(1)
	}

	repository := helpers.ClaudeStorage(helpers.DetectClaudeDevice())

	listCmd := exec.Command("restic", "snapshots", "--repo", repository, "--insecure-no-password", "--tag", "claude", "--json")
	listCmd.Stderr = os.Stderr
	listOut, err := listCmd.Output()
	if err != nil {
		fmt.Println(aurora.Red("Error: Failed to list snapshots: " + err.Error()))
		os.Exit(1)
	}

	var snapshots []snapshot
	if err := json.Unmarshal(listOut, &snapshots); err != nil {
		fmt.Println(aurora.Red("Error: Failed to parse snapshots: " + err.Error()))
		os.Exit(1)
	}

	if len(snapshots) == 0 {
		fmt.Println(aurora.Red("Error: No pushes found in " + repository))
		os.Exit(1)
	}

	slices.SortFunc(snapshots, func(a, b snapshot) int {
		return a.Time.Compare(b.Time)
	})

	cacheDir := helpers.ResolvePath("~/.cache/agents-stats")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		panic(err)
	}

	rawDir := filepath.Join(cacheDir, "raw")
	restoredFile := filepath.Join(cacheDir, ".restored")
	restoredIDs := map[string]bool{}
	if utils.IsFileExists(restoredFile) {
		restoredData, err := os.ReadFile(restoredFile)
		if err != nil {
			panic(err)
		}
		for _, id := range strings.Fields(string(restoredData)) {
			restoredIDs[id] = true
		}
	}

	for _, snap := range snapshots {
		if restoredIDs[snap.ID] {
			continue
		}

		helpers.ExecNativeCommand(
			[]string{"restic", "restore", "--repo", repository, "--insecure-no-password", snap.ID + ":/projects", "--target", filepath.Join(rawDir, snap.Hostname), "--include", "*.jsonl", "--overwrite", "if-newer", "--quiet"},
			helpers.ExecCommandOptions{Exit: true},
		)

		file, err := os.OpenFile(restoredFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			panic(err)
		}
		if _, err := file.WriteString(snap.ID + "\n"); err != nil {
			panic(err)
		}
		file.Close()
	}

	responses := map[string]*response{}
	var unreadableLines int

	err = filepath.WalkDir(rawDir, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if dirEntry.IsDir() || filepath.Ext(path) != ".jsonl" {
			return nil
		}

		relPath, err := filepath.Rel(rawDir, path)
		if err != nil {
			return err
		}
		device, _, found := strings.Cut(relPath, string(filepath.Separator))
		if !found {
			return fmt.Errorf("unexpected file outside a device folder: %s", path)
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		reader := bufio.NewReaderSize(file, 1024*1024)
		for lineNumber := 1; ; lineNumber++ {
			line, readErr := reader.ReadBytes('\n')
			if len(line) > 0 && bytes.Contains(line, []byte(`"assistant"`)) && bytes.Contains(line, []byte(`"usage"`)) {
				var parsed entry
				if err := json.Unmarshal(line, &parsed); err != nil {
					unreadableLines++
				} else if parsed.Type == "assistant" && parsed.Message.Usage != nil && parsed.Message.Model != "" && parsed.Message.Model != "<synthetic>" {
					usage := parsed.Message.Usage
					current := tokens{
						Requests:   1,
						Input:      usage.InputTokens,
						Output:     usage.OutputTokens,
						CacheWrite: usage.CacheCreationInputTokens,
						CacheRead:  usage.CacheReadInputTokens,
					}

					var key string
					if parsed.Message.ID != "" {
						key = parsed.Message.ID + ":" + parsed.RequestID
					} else {
						key = path + ":" + strconv.Itoa(lineNumber)
					}

					if existing, ok := responses[key]; ok {
						existing.Tokens.Input = max(existing.Tokens.Input, current.Input)
						existing.Tokens.Output = max(existing.Tokens.Output, current.Output)
						existing.Tokens.CacheWrite = max(existing.Tokens.CacheWrite, current.CacheWrite)
						existing.Tokens.CacheRead = max(existing.Tokens.CacheRead, current.CacheRead)
					} else {
						month := "unknown"
						if timestamp, err := time.Parse(time.RFC3339, parsed.Timestamp); err == nil {
							month = timestamp.Local().Format("2006-01")
						}
						responses[key] = &response{Device: device, Model: parsed.Message.Model, Month: month, Tokens: current}
					}
				}
			}

			if errors.Is(readErr, io.EOF) {
				return nil
			}
			if readErr != nil {
				return readErr
			}
		}
	})
	if err != nil {
		panic(err)
	}

	var total tokens
	byDevice := map[string]*tokens{}
	byModel := map[string]*tokens{}
	byMonth := map[string]*tokens{}
	for _, resp := range responses {
		addTokens(&total, resp.Tokens)
		addTokens(groupTokens(byDevice, resp.Device), resp.Tokens)
		addTokens(groupTokens(byModel, resp.Model), resp.Tokens)
		addTokens(groupTokens(byMonth, resp.Month), resp.Tokens)
	}

	monthNames := slices.Sorted(maps.Keys(byMonth))
	data := report{
		Storage:         repository,
		Generated:       time.Now().Format("2006-01-02 15:04"),
		Pushes:          len(snapshots),
		UnreadableLines: unreadableLines,
		Total:           makeRow("Total", total, total),
		Sections: []section{
			{Title: "Devices", Label: "Device", Rows: makeRows(byDevice, sortedByTotal(byDevice), total)},
			{Title: "Models", Label: "Model", Rows: makeRows(byModel, sortedByTotal(byModel), total)},
			{Title: "Months", Label: "Month", Rows: makeRows(byMonth, monthNames, total)},
		},
	}

	reportPath := filepath.Join(os.TempDir(), "agents-stats.html")
	reportFile, err := os.Create(reportPath)
	if err != nil {
		panic(err)
	}
	if err := reportTemplate.Execute(reportFile, data); err != nil {
		panic(err)
	}
	if err := reportFile.Close(); err != nil {
		panic(err)
	}

	fmt.Println(aurora.Green("Claude stats written to " + reportPath))
	helpers.Open(reportPath, helpers.ExecCommandOptions{Exit: true})
}

type row struct {
	Name       string
	Requests   string
	Input      string
	Output     string
	CacheWrite string
	CacheRead  string
	Total      string
	Share      string
}

type section struct {
	Title string
	Label string
	Rows  []row
}

type report struct {
	Storage         string
	Generated       string
	Pushes          int
	UnreadableLines int
	Total           row
	Sections        []section
}

func groupTokens(groups map[string]*tokens, name string) *tokens {
	if groups[name] == nil {
		groups[name] = &tokens{}
	}
	return groups[name]
}

func addTokens(target *tokens, source tokens) {
	target.Requests += source.Requests
	target.Input += source.Input
	target.Output += source.Output
	target.CacheWrite += source.CacheWrite
	target.CacheRead += source.CacheRead
}

func sumTokens(values tokens) int64 {
	return values.Input + values.Output + values.CacheWrite + values.CacheRead
}

func sortedByTotal(groups map[string]*tokens) []string {
	names := slices.Collect(maps.Keys(groups))
	slices.SortFunc(names, func(a, b string) int {
		return cmp.Compare(sumTokens(*groups[b]), sumTokens(*groups[a]))
	})
	return names
}

func makeRows(groups map[string]*tokens, names []string, total tokens) []row {
	var rows []row
	for _, name := range names {
		rows = append(rows, makeRow(name, *groups[name], total))
	}
	return rows
}

func makeRow(name string, values tokens, total tokens) row {
	share := 0.0
	if sumTokens(total) > 0 {
		share = float64(sumTokens(values)) / float64(sumTokens(total)) * 100
	}

	return row{
		Name:       name,
		Requests:   formatCount(values.Requests),
		Input:      formatTokens(values.Input),
		Output:     formatTokens(values.Output),
		CacheWrite: formatTokens(values.CacheWrite),
		CacheRead:  formatTokens(values.CacheRead),
		Total:      formatTokens(sumTokens(values)),
		Share:      strconv.FormatFloat(share, 'f', 1, 64),
	}
}

func formatCount(count int64) string {
	digits := strconv.FormatInt(count, 10)
	var result strings.Builder
	for i, digit := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			result.WriteByte(',')
		}
		result.WriteRune(digit)
	}
	return result.String()
}

func formatTokens(count int64) string {
	switch {
	case count >= 1_000_000_000:
		return fmt.Sprintf("%.2fB", float64(count)/1_000_000_000)
	case count >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(count)/1_000_000)
	case count >= 1_000:
		return fmt.Sprintf("%.1fK", float64(count)/1_000)
	default:
		return strconv.FormatInt(count, 10)
	}
}

var reportTemplate = template.Must(template.New("report").Parse(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Claude stats</title>
<style>
  :root {
    color-scheme: light dark;
    --bg: #fbfaf8;
    --text: #1c1b19;
    --muted: #77736c;
    --line: #e7e3dc;
    --bar: #d97757;
    --bar-track: #f0ebe4;
  }
  @media (prefers-color-scheme: dark) {
    :root {
      --bg: #171614;
      --text: #ece9e4;
      --muted: #8f8a82;
      --line: #2c2a27;
      --bar: #e08a6b;
      --bar-track: #262421;
    }
  }
  * { box-sizing: border-box; }
  body {
    margin: 0;
    background: var(--bg);
    color: var(--text);
    font: 15px/1.5 -apple-system, BlinkMacSystemFont, "Segoe UI", system-ui, sans-serif;
  }
  main { max-width: 1040px; margin: 0 auto; padding: 56px 32px 80px; }
  h1 { font-size: 15px; font-weight: 600; margin: 0; color: var(--muted); }
  .total { font-size: 72px; font-weight: 650; letter-spacing: -0.03em; line-height: 1.05; margin: 8px 0 4px; font-variant-numeric: tabular-nums; }
  .meta { color: var(--muted); margin: 0; }
  .meta span { display: inline-block; white-space: nowrap; }
  .meta span + span::before { content: "·"; margin: 0 8px; }
  .breakdown { display: flex; flex-wrap: wrap; gap: 8px 40px; margin: 32px 0 0; padding: 20px 0; border-top: 1px solid var(--line); border-bottom: 1px solid var(--line); }
  .breakdown div { display: flex; flex-direction: column; }
  .breakdown dt { color: var(--muted); font-size: 13px; }
  .breakdown dd { margin: 0; font-size: 20px; font-weight: 600; font-variant-numeric: tabular-nums; }
  .warning { margin: 16px 0 0; color: #b45309; }
  h2 { font-size: 18px; font-weight: 600; margin: 56px 0 12px; }
  .scroll { overflow-x: auto; }
  table { width: 100%; min-width: 860px; table-layout: fixed; border-collapse: collapse; font-variant-numeric: tabular-nums; }
  th:first-child { width: 22%; }
  th:last-child { width: 20%; }
  th { font-size: 12px; font-weight: 500; color: var(--muted); text-align: right; padding: 8px 12px; border-bottom: 1px solid var(--line); white-space: nowrap; }
  td { text-align: right; padding: 10px 12px; border-bottom: 1px solid var(--line); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  th:first-child, td:first-child { text-align: left; padding-left: 0; }
  th:last-child, td:last-child { padding-right: 0; }
  td.name { font-weight: 500; }
  td.strong { font-weight: 600; }
  .share { display: flex; align-items: center; justify-content: flex-end; gap: 10px; min-width: 160px; }
  .share .track { flex: 1; max-width: 120px; height: 6px; border-radius: 3px; background: var(--bar-track); overflow: hidden; }
  .share .fill { height: 100%; background: var(--bar); border-radius: 3px; }
  .share .value { width: 48px; color: var(--muted); }
  @media (max-width: 640px) {
    main { padding: 32px 16px 56px; }
    .total { font-size: 52px; }
  }
</style>
</head>
<body>
<main>
  <h1>Claude tokens</h1>
  <p class="total">{{.Total.Total}}</p>
  <p class="meta"><span>{{.Total.Requests}} requests</span><span>{{.Pushes}} pushes</span><span>{{.Storage}}</span><span>Generated {{.Generated}}</span></p>

  <dl class="breakdown">
    <div><dt>Input</dt><dd>{{.Total.Input}}</dd></div>
    <div><dt>Output</dt><dd>{{.Total.Output}}</dd></div>
    <div><dt>Cache write</dt><dd>{{.Total.CacheWrite}}</dd></div>
    <div><dt>Cache read</dt><dd>{{.Total.CacheRead}}</dd></div>
  </dl>
  {{if .UnreadableLines}}<p class="warning">Skipped {{.UnreadableLines}} unreadable lines.</p>{{end}}

  {{range .Sections}}
  <h2>{{.Title}}</h2>
  <div class="scroll">
    <table>
      <thead><tr><th>{{.Label}}</th><th>Requests</th><th>Input</th><th>Output</th><th>Cache write</th><th>Cache read</th><th>Total</th><th>Share</th></tr></thead>
      <tbody>
      {{range .Rows}}
        <tr>
          <td class="name">{{.Name}}</td>
          <td>{{.Requests}}</td>
          <td>{{.Input}}</td>
          <td>{{.Output}}</td>
          <td>{{.CacheWrite}}</td>
          <td>{{.CacheRead}}</td>
          <td class="strong">{{.Total}}</td>
          <td><div class="share"><div class="track"><div class="fill" style="width: {{.Share}}%"></div></div><span class="value">{{.Share}}%</span></div></td>
        </tr>
      {{end}}
      </tbody>
    </table>
  </div>
  {{end}}
</main>
</body>
</html>
`))
