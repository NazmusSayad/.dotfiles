package main

import (
	"cmp"
	"fmt"
	"html/template"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"dotfiles/src/helpers"
	"dotfiles/src/helpers/agentstore"

	"github.com/logrusorgru/aurora/v4"
)

type tokens struct {
	Requests   int64
	Input      int64
	Output     int64
	CacheWrite int64
	CacheRead  int64
}

func main() {
	stats := agentstore.FetchStats(agentstore.LoadConfig(os.Args[1:]))

	var total tokens
	byDevice := map[string]*tokens{}
	byModel := map[string]*tokens{}
	byMonth := map[string]*tokens{}
	for _, group := range stats.Groups {
		current := tokens{
			Requests:   group.Requests,
			Input:      group.Input,
			Output:     group.Output,
			CacheWrite: group.CacheWrite,
			CacheRead:  group.CacheRead,
		}
		addTokens(&total, current)
		addTokens(groupTokens(byDevice, group.Device), current)
		addTokens(groupTokens(byModel, group.Model), current)
		addTokens(groupTokens(byMonth, group.Month), current)
	}

	monthNames := slices.Sorted(maps.Keys(byMonth))
	data := report{
		Storage:         stats.Storage,
		Generated:       time.Now().Format("2006-01-02 15:04"),
		Sessions:        formatCount(int64(stats.Sessions)),
		UnreadableLines: stats.UnreadableLines,
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
	Sessions        string
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
  <p class="meta"><span>{{.Total.Requests}} requests</span><span>{{.Sessions}} sessions</span><span>{{.Storage}}</span><span>Generated {{.Generated}}</span></p>

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
