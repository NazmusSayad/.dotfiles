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
	if ok && existingValue == value {
		fmt.Println(aurora.Blue(name).String(), aurora.Green("already set"))
		return
	}

	helpers.WriteEnv(helpers.ScopeUser, name, value)
}

func main() {
	initGoEnv()
	initJavaEnv()
	initAndroidSdkEnv()
}

func initGoEnv() {
	goBinCmd := exec.Command("go", "env", "GOBIN")
	goPathCmd := exec.Command("go", "env", "GOPATH")
	goRootCmd := exec.Command("go", "env", "GOROOT")

	goBinOutput, err := goBinCmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		setEnv("GOBIN", strings.TrimSpace(string(goBinOutput)))
	}

	goPathOutput, err := goPathCmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		setEnv("GOPATH", strings.TrimSpace(string(goPathOutput)))
	}

	goRootOutput, err := goRootCmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		setEnv("GOROOT", strings.TrimSpace(string(goRootOutput)))
	}
}

func initJavaEnv() {
	javaHomeCmd := exec.Command("mise", "where", "java")
	javaHomeOutput, err := javaHomeCmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	setEnv("JAVA_HOME", strings.TrimSpace(string(javaHomeOutput)))
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
