package main

import (
	"dotfiles/src/helpers"
	"dotfiles/src/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

type MiseEnv struct {
	Source string `json:"source"`
	Tool   string `json:"tool"`
	Value  string `json:"value"`
}

func setEnv(name, value string) {
	fmt.Println(aurora.Blue(name).String(), aurora.Green(value))

	existingValue, ok := os.LookupEnv(name)
	if ok && existingValue != "" && existingValue == value {
		fmt.Println(aurora.Blue(name).String(), aurora.Green("already set"))
		return
	}

	helpers.WriteEnv(helpers.ScopeUser, name, value)
}

func main() {
	initMiseEnv()
	initAndroidSdkEnv()
}

func initMiseEnv() {
	miseEnvCmd := exec.Command("mise", "env", "--dotenv")
	miseEnvOutput, err := miseEnvCmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	miseEnvLines := strings.Split(strings.TrimSpace(string(miseEnvOutput)), "\n")
	for _, line := range miseEnvLines {
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			fmt.Println("Error:", line)
			continue
		}

		setEnv(strings.TrimSpace(key), strings.TrimSpace(value))
	}
}

func initAndroidSdkEnv() {
	androidSdkPath := filepath.Join(os.Getenv("LOCALAPPDATA"), "Android", "Sdk")
	if !utils.IsFileExists(androidSdkPath) {
		fmt.Println("Android SDK not found")
		return
	}

	setEnv("ANDROID_HOME", androidSdkPath)
	setEnv("ANDROID_SDK_ROOT", androidSdkPath)
}
