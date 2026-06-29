package helpers

import (
	"os"
	"strings"
)

func GetBrewTaps(brewFilePath string) []string {
	brewFile, _ := os.ReadFile(brewFilePath)
	brewFileLines := strings.Split(string(brewFile), "\n")

	brewFileTaps := []string{}
	for _, line := range brewFileLines {
		if after, ok := strings.CutPrefix(line, "tap"); ok {
			tap := strings.TrimSpace(after)
			tap = strings.Trim(tap, "\"")
			brewFileTaps = append(brewFileTaps, tap)
		}
	}

	return brewFileTaps
}
