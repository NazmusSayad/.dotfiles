package main

import (
	helpers "dotfiles/src"
	"dotfiles/src/winget"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
)

func main() {
	packages := winget.GetWingetPackages("./config/winget-apps.jsonc")
	updatablePackages := getUpdatablePackageIDs()

	for _, p := range packages {
		if p.IgnoreUpgrade || p.Version != "" {
			fmt.Println("\n- Skipping", p.ID)
			continue
		}

		if !slices.Contains(updatablePackages, p.ID) {
			fmt.Println("\n- No updates available for", p.ID)
			continue
		}

		fmt.Println("\n- Upgrading", p.ID)
		parts := winget.BuildWingetUpgradeCommands(p)
		cmd := exec.Command(parts[0], parts[1:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
	}

	fmt.Println("\nDone!")
	helpers.PressAnyKeyOrWaitToExit()
}

func getUpdatablePackageIDs() []string {
	cmd := exec.Command("winget", "upgrade")
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	stringifiedOutput := string(output)
	println(stringifiedOutput)

	lines := strings.Split(stringifiedOutput, "\n")
	var packageIDs []string

	foundSeparator := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "----") {
			foundSeparator = true
			continue
		}
		if !foundSeparator || line == "" || !strings.Contains(line, "winget") {
			continue
		}

		wingetIndex := strings.LastIndex(line, "winget")
		if wingetIndex >= 0 {
			line = strings.TrimSpace(line[:wingetIndex])
		}

		re := regexp.MustCompile(`\b([a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+)\b`)
		matches := re.FindAllString(line, -1)

		for _, match := range matches {
			parts := strings.Split(match, ".")
			if len(parts) == 2 && !strings.Contains(parts[0], "-") && !strings.Contains(parts[1], "-") {
				firstChar := parts[0][0]
				if len(parts[0]) > 1 && len(parts[1]) > 1 && (firstChar < '0' || firstChar > '9') {
					packageIDs = append(packageIDs, match)
					break
				}
			}
		}
	}

	return packageIDs
}
