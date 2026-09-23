package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"dotfiles/src/helpers"
	"dotfiles/src/utils"

	"github.com/joho/godotenv"
	"github.com/logrusorgru/aurora/v4"
)

func main() {
	initOS()
	initMiseEnv()
	initAndroidSdkEnv()

	if err := initDotEnv(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func initOS() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	setEnv("HOME", homeDir)

	bashPath, err := exec.LookPath("bash")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	setEnv("SHELL", bashPath)
}

func initMiseEnv() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	miseEnvCmd := exec.Command("mise", "env", "--dotenv", "--cd", homeDir)
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

func initDotEnv() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dotfilesEnvPath := helpers.ResolvePath("@/.env")
	env := make(map[string]string)
	for _, path := range []string{dotfilesEnvPath, filepath.Join(homeDir, ".env")} {
		values, err := godotenv.Read(path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		for name, value := range values {
			env[strings.ToUpper(name)] = value
		}
	}

	lockPath := dotfilesEnvPath + ".lock"
	lockData, err := os.ReadFile(lockPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("read %s: %w", lockPath, err)
	}

	tracked := make(map[string]bool)
	for _, name := range strings.Split(string(lockData), "\n") {
		name = strings.ToUpper(strings.TrimSpace(name))
		if name != "" {
			tracked[name] = true
		}
	}

	lockFile, err := os.CreateTemp(filepath.Dir(lockPath), ".env.lock-*")
	if err != nil {
		return fmt.Errorf("create lock file: %w", err)
	}
	defer os.Remove(lockFile.Name())
	defer lockFile.Close()

	var syncErr error
	for name, value := range env {
		if _, err := helpers.WriteEnv(helpers.ScopeUser, name, value); err != nil {
			syncErr = errors.Join(syncErr, fmt.Errorf("set %s: %w", name, err))
			continue
		}

		tracked[name] = true
		fmt.Println(aurora.Blue(name).String(), aurora.Green("set from .env"))
	}

	for name := range tracked {
		if _, ok := env[name]; ok {
			continue
		}
		if _, err := helpers.WriteEnv(helpers.ScopeUser, name, ""); err != nil {
			syncErr = errors.Join(syncErr, fmt.Errorf("delete %s: %w", name, err))
			continue
		}

		delete(tracked, name)
		fmt.Println(aurora.Blue(name).String(), aurora.Yellow("removed from user environment"))
	}

	names := make([]string, 0, len(tracked))
	for name := range tracked {
		names = append(names, name)
	}
	sort.Strings(names)

	if _, err := lockFile.WriteString(strings.Join(names, "\n")); err != nil {
		return errors.Join(syncErr, fmt.Errorf("write %s: %w", lockPath, err))
	}
	if err := lockFile.Close(); err != nil {
		return errors.Join(syncErr, fmt.Errorf("close %s: %w", lockPath, err))
	}
	if err := os.Rename(lockFile.Name(), lockPath); err != nil {
		return errors.Join(syncErr, fmt.Errorf("replace %s: %w", lockPath, err))
	}

	return syncErr
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
