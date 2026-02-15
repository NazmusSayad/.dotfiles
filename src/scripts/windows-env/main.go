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
	initGoEnv()
	initJavaEnv()
	initAndroidSdkEnv()
}    

func initGoEnv() {
	goPath := getGoEnv("GOPATH", os.Environ())
	goEnv := getEnvWithoutGoVars() 
   
	goBin := getGoEnv("GOBIN", goEnv)
	goRoot := getGoEnv("GOROOT", goEnv)

	setEnv("GOPATH", goPath) 
	setEnv("GOBIN", goBin)
	setEnv("GOROOT", goRoot)
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

func getGoEnv(name string, env []string) string {
	cmd := exec.Command("go", "env", name)
	cmd.Env = env

	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(output))
}

func getEnvWithoutGoVars() []string {
	env := os.Environ()
	filtered := make([]string, 0, len(env))

	for _, item := range env {
		key, _, ok := strings.Cut(item, "=")
		if !ok || strings.HasPrefix(key, "GO") {
			continue
		}

		filtered = append(filtered, item)
	}

	return filtered
}
