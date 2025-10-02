package winget

import (
	"os/exec"
	"regexp"
	"strings"
)

type UpdatablePackage struct {
	ID        string
	Version   string
	Available string
}

func GetUpdatablePackageIDs() []string {
	cmd := exec.Command("winget", "upgrade")
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	lines := strings.Split(string(output), "\n")
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
