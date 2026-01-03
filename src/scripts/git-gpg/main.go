package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	if _, gitPathError := exec.LookPath("git"); gitPathError != nil {
		fmt.Println(aurora.Red("Error: Git not installed"))
		os.Exit(1)
	}

	gpgPath, gpgPathError := exec.LookPath("gpg")
	if gpgPathError != nil {
		fmt.Println(aurora.Red("Error: GPG not installed"))
		os.Exit(1)
	}

	gitNameOut, _ := exec.Command("git", "config", "--get", "user.name").Output()
	gitName := strings.TrimSpace(string(gitNameOut))

	gitEmailOut, _ := exec.Command("git", "config", "--get", "user.email").Output()
	gitEmail := strings.TrimSpace(string(gitEmailOut))

	if gitEmail == "" || gitName == "" {
		fmt.Println(aurora.Red("Error: Git user.email or user.name not configured"))
		fmt.Println("Please run: git config --global user.name \"Your Name\"")
		fmt.Println("Please run: git config --global user.email \"your@email.com\"")
		os.Exit(1)
	}

	fmt.Println("User name      :", gitName)
	fmt.Println("User email     :", gitEmail)

	listKeysOut, _ := exec.Command("gpg", "--list-secret-keys", "--keyid-format", "LONG").Output()
	hasKeys := strings.Contains(string(listKeysOut), "sec")

	if !hasKeys {
		fmt.Println(aurora.Yellow("No GPG keys found, generating new key..."))

		batchFilePath := filepath.Join(os.TempDir(), "gpg_batch.txt")
		batchContent := strings.Join([]string{
			"Key-Type: RSA",
			"Key-Length: 4096",
			"Key-Usage: sign",
			"Name-Real: " + gitName,
			"Name-Email: " + gitEmail,
			"Expire-Date: 0",
			"%no-protection",
			"%commit",
			"",
		}, "\n")

		fileWriteErr := os.WriteFile(batchFilePath, []byte(batchContent), 0644)
		if fileWriteErr != nil {
			fmt.Println(aurora.Red("Error: failed to write batch file: " + fileWriteErr.Error()))
			os.Exit(1)
		}

		generateCmd := exec.Command("gpg", "--batch", "--generate-key", batchFilePath)
		generateCmdErr := generateCmd.Run()
		os.Remove(batchFilePath)

		if generateCmdErr != nil {
			os.Exit(1)
		}
	}

	listKeysOut, _ = exec.Command("gpg", "--list-secret-keys", "--keyid-format", "LONG").Output()

	var gpgKeyID string
	for line := range strings.SplitSeq(string(listKeysOut), "\n") {
		if strings.Contains(line, "sec") {
			parts := strings.Split(line, "/")
			if len(parts) > 1 {
				keyPart := strings.Fields(parts[1])[0]
				gpgKeyID = keyPart
				break
			}
		}
	}

	if gpgKeyID == "" {
		os.Exit(1)
	}

	exec.Command("git", "config", "--global", "user.signingkey", gpgKeyID).Run()
	exec.Command("git", "config", "--global", "commit.gpgsign", "true").Run()
	exec.Command("git", "config", "--global", "gpg.program", gpgPath).Run()

	exportCmd := exec.Command("gpg", "--armor", "--export", gpgKeyID)
	exportCmd.Stdout = os.Stdout
	exportCmd.Stderr = os.Stderr
	exportCmd.Run()
}
