package helpers

import (
	"archive/zip"
	"bytes"
	"dotfiles/src/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/logrusorgru/aurora/v4"
	"gopkg.in/yaml.v3"
)

type Asset struct {
	Name               string `json:"name"`
	Size               int64  `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
	ContentType        string `json:"content_type"`
}

type Release struct {
	TagName     string  `json:"tag_name"`
	Name        string  `json:"name"`
	PublishedAt string  `json:"published_at"`
	URL         string  `json:"html_url"`
	Assets      []Asset `json:"assets"`
}

func getLatestReleaseZipURL(repoURL, pattern string) (Asset, error) {
	parts := strings.Split(strings.TrimSuffix(strings.TrimPrefix(repoURL, "https://github.com/"), "/"), "/")
	if len(parts) < 2 {
		return Asset{}, fmt.Errorf("invalid GitHub URL format")
	}
	owner, repo := parts[0], parts[1]

	re, err := regexp.Compile(pattern)
	if err != nil {
		return Asset{}, fmt.Errorf("invalid regex pattern: %v", err)
	}

	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	resp, err := http.Get(apiURL)
	if err != nil {
		return Asset{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Asset{}, fmt.Errorf("received status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Asset{}, err
	}

	var release Release
	if err := json.Unmarshal(body, &release); err != nil {
		return Asset{}, err
	}

	for _, asset := range release.Assets {
		if re.MatchString(asset.Name) {
			fmt.Println(aurora.Faint("Found file: " + asset.Name))
			time.Sleep(3 * time.Second)
			return asset, nil
		}
	}

	return Asset{}, fmt.Errorf("file matching pattern '%s' not found", pattern)
}

func downloadGithubReleaseFile(outDir, ghURL, pattern string) (Asset, []byte, error) {
	asset, err := getLatestReleaseZipURL(ghURL, pattern)
	if err != nil {
		return Asset{}, nil, err
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return Asset{}, nil, err
	}

	resp, err := http.Get(asset.BrowserDownloadURL)
	if err != nil {
		return Asset{}, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Asset{}, nil, fmt.Errorf("received status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Asset{}, nil, err
	}

	return asset, body, nil
}

func WriteGithubReleaseFile(outDir, ghURL, pattern string) error {
	outputPath := filepath.Join(outDir, filepath.Base(pattern))
	if utils.IsFileExists(outputPath) {
		fmt.Println(aurora.Faint("File already exists: " + outputPath))
		return nil
	} else {
		fmt.Println(aurora.Faint("Downloading file: " + outputPath))
	}

	asset, body, err := downloadGithubReleaseFile(outDir, ghURL, pattern)
	if err != nil {
		fmt.Println(aurora.Red("Error: " + err.Error()))
		os.Exit(1)
	}

	filename := filepath.Base(asset.Name)
	return os.WriteFile(filepath.Join(outDir, filename), body, 0644)
}

func WriteGithubReleaseZipFile(outDir, ghURL, archivePattern, exeName string) error {
	outputPath := filepath.Join(outDir, exeName)
	if utils.IsFileExists(outputPath) {
		fmt.Println(aurora.Faint("File already exists: " + outputPath))
		return nil
	} else {
		fmt.Println(aurora.Faint("Downloading file: " + outputPath))
	}

	_, body, err := downloadGithubReleaseFile(outDir, ghURL, archivePattern)
	if err != nil {
		fmt.Println(aurora.Red("Error: failed to download file: " + err.Error()))
		os.Exit(1)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return err
	}

	var targetFile *zip.File
	for _, f := range zipReader.File {
		if filepath.Base(f.Name) == exeName {
			targetFile = f
			break
		}
	}

	if targetFile == nil {
		return fmt.Errorf("file '%s' not found in zip", exeName)
	}

	rc, err := targetFile.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, rc)
	return err
}

type GhHostConfig struct {
	User string `yaml:"user"`
}

func GetGitHubUser() string {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return ""
	}

	hostsPath := filepath.Join(appData, "GitHub CLI", "hosts.yml")
	if !utils.IsFileExists(hostsPath) {
		return ""
	}

	data, err := os.ReadFile(hostsPath)
	if err != nil {
		return ""
	}

	var hosts map[string]GhHostConfig
	if err := yaml.Unmarshal(data, &hosts); err != nil {
		return ""
	}

	if config, ok := hosts["github.com"]; ok {
		return config.User
	}

	return ""
}

func GetGitHubUserOrExit() string {
	user := GetGitHubUser()
	if user == "" {
		fmt.Println(aurora.Red("No GitHub user found"))
		os.Exit(1)
	}

	return user
}
