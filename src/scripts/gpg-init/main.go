package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"dotfiles/src/utils"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	nameFlag := flag.String("name", "", "Key user name (default: git config user.name)")
	emailFlag := flag.String("email", "", "Key user email (default: git config user.email)")
	fileFlag := flag.String("file", "", "Git config file to write (default: global config)")
	flag.Parse()

	if !utils.IsCommandInPath("git") {
		fmt.Println(aurora.Red("Error: Git not installed"))
		os.Exit(1)
	}

	if !utils.IsCommandInPath("gpg") {
		fmt.Println(aurora.Red("Error: GPG not installed"))
		os.Exit(1)
	}

	gitName := *nameFlag
	if gitName == "" {
		gitNameOut, _ := exec.Command("git", "config", "--get", "user.name").Output()
		gitName = strings.TrimSpace(string(gitNameOut))
	}

	gitEmail := *emailFlag
	if gitEmail == "" {
		gitEmailOut, _ := exec.Command("git", "config", "--get", "user.email").Output()
		gitEmail = strings.TrimSpace(string(gitEmailOut))
	}

	if gitEmail == "" || gitName == "" {
		fmt.Println(aurora.Red("Error: Git user.email or user.name not configured"))
		fmt.Println("Please run: git config --global user.name \"Your Name\"")
		fmt.Println("Please run: git config --global user.email \"your@email.com\"")
		os.Exit(1)
	}

	fmt.Println("User name      :", gitName)
	fmt.Println("User email     :", gitEmail)

	keyQuery := "<" + gitEmail + ">"

	listKeysOut, _ := exec.Command("gpg", "--list-secret-keys", "--keyid-format", "LONG", keyQuery).Output()
	hasKeys := strings.Contains(string(listKeysOut), "sec ")

	if !hasKeys {
		fmt.Println(aurora.Yellow(">> No GPG keys found, generating new key..."))

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

		generateCmd := exec.Command("gpg", "--batch", "--generate-key")
		generateCmd.Stdin = strings.NewReader(batchContent)
		generateCmd.Stdout = os.Stdout
		generateCmd.Stderr = os.Stderr
		if err := generateCmd.Run(); err != nil {
			os.Exit(1)
		}
	}

	listKeysOut, _ = exec.Command("gpg", "--list-secret-keys", "--keyid-format", "LONG", keyQuery).Output()

	var gpgKeyID string
	for line := range strings.SplitSeq(string(listKeysOut), "\n") {
		if strings.HasPrefix(line, "sec ") {
			parts := strings.Split(line, "/")
			if len(parts) > 1 {
				gpgKeyID = strings.Fields(parts[1])[0]
				break
			}
		}
	}

	if gpgKeyID == "" {
		os.Exit(1)
	}

	configScope := []string{"config", "--global"}
	if *fileFlag != "" {
		configScope = []string{"config", "--file", *fileFlag}
	}

	exec.Command("git", append(configScope, "user.signingkey", gpgKeyID)...).Run()
	exec.Command("git", append(configScope, "gpg.program", "gpg")...).Run()
	exec.Command("git", append(configScope, "gpg.format", "openpgp")...).Run()

	exec.Command("git", append(configScope, "commit.gpgsign", "true")...).Run()
	exec.Command("git", append(configScope, "tag.gpgsign", "true")...).Run()

	exportCmd := exec.Command("gpg", "--armor", "--export", gpgKeyID)
	exportCmd.Stdout = os.Stdout
	exportCmd.Stderr = os.Stderr

	fmt.Println("")
	exportCmd.Run()
}
