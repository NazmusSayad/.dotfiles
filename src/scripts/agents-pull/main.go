package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
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

	slices.SortFunc(snapshots, func(a, b snapshot) int {
		return a.Time.Compare(b.Time)
	})

	archiveDir := helpers.ResolvePath("~/.agents-archive")
	if err := os.MkdirAll(archiveDir, 0o755); err != nil {
		panic(err)
	}

	pulledFile := filepath.Join(archiveDir, ".pulled")
	pulledIDs := map[string]bool{}
	if utils.IsFileExists(pulledFile) {
		pulledData, err := os.ReadFile(pulledFile)
		if err != nil {
			panic(err)
		}
		for _, id := range strings.Fields(string(pulledData)) {
			pulledIDs[id] = true
		}
	}

	for _, snap := range snapshots {
		if pulledIDs[snap.ID] {
			continue
		}

		target := filepath.Join(archiveDir, snap.Hostname)
		fmt.Println(aurora.Blue(snap.Hostname), aurora.Faint(snap.Time.Format(time.DateTime)))

		helpers.ExecNativeCommand(
			[]string{"restic", "restore", "--repo", repository, "--insecure-no-password", snap.ID, "--target", target, "--overwrite", "if-newer"},
			helpers.ExecCommandOptions{Exit: true},
		)

		file, err := os.OpenFile(pulledFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			panic(err)
		}
		if _, err := file.WriteString(snap.ID + "\n"); err != nil {
			panic(err)
		}
		file.Close()
	}

	fmt.Println(aurora.Green("Claude data pulled into " + archiveDir))
}
