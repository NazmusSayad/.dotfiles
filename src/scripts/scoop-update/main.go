package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"dotfiles/src/constants"
	"dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
)

type scoopVersionConfig struct {
	URL   string `yaml:"url"`
	Regex string `yaml:"regex"`
	Fixed string `yaml:"fixed"`
}

type scoopAutoupdateArchConfig struct {
	URL string `yaml:"url" json:"url"`
}

type scoopAutoupdateConfig struct {
	Architecture map[string]scoopAutoupdateArchConfig `yaml:"architecture" json:"architecture"`
}

type scoopAppTemplate struct {
	Description string                `yaml:"description"`
	Homepage    string                `yaml:"homepage"`
	Extract     string                `yaml:"extract"`
	Install     []string              `yaml:"install"`
	Bin         [][]string            `yaml:"bin"`
	Shortcuts   [][]string            `yaml:"shortcuts"`
	Version     scoopVersionConfig    `yaml:"version"`
	Autoupdate  scoopAutoupdateConfig `yaml:"autoupdate"`
}

type githubReleaseAsset struct {
	BrowserDownloadURL string `json:"browser_download_url"`
	Digest             string `json:"digest"`
}

type githubRelease struct {
	TagName string               `json:"tag_name"`
	Assets  []githubReleaseAsset `json:"assets"`
}

func main() {
	apps := helpers.ReadConfig[map[string]scoopAppTemplate]("@/config/scoop-apps.yaml")
	resolvedScoopDir := helpers.ResolvePath("@/" + constants.SCOOP_DIR)
	githubToken := strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	fmt.Println("Updating Scoop manifests in", aurora.Cyan(resolvedScoopDir))

	err := os.MkdirAll(resolvedScoopDir, 0o755)
	if err != nil {
		fmt.Println(aurora.Red("Failed to create Scoop directory:"), err)
		os.Exit(1)
	}

	successCount := 0
	failCount := 0

	for appID, app := range apps {
		if app.Version.URL == "" || (app.Version.Regex == "" && app.Version.Fixed == "") {
			fmt.Println(aurora.Red("Failed:"), appID, "missing version.url and version.regex/version.fixed")
			failCount++
			continue
		}

		var versionRegex *regexp.Regexp
		if app.Version.Regex != "" {
			compiledRegex, compileErr := regexp.Compile(app.Version.Regex)
			if compileErr != nil {
				fmt.Println(aurora.Red("Failed:"), appID, "invalid version regex:", compileErr)
				failCount++
				continue
			}

			versionRegex = compiledRegex
		}

		version := app.Version.Fixed

		parsedVersionURL, parseErr := url.Parse(strings.TrimSpace(app.Version.URL))
		if parseErr != nil {
			fmt.Println(aurora.Red("Failed:"), appID, "invalid version.url:", parseErr)
			failCount++
			continue
		}

		parsedVersionURL.Path = strings.TrimSuffix(parsedVersionURL.Path, "/")
		versionURL := parsedVersionURL.String()
		if app.Version.Fixed != "" {
			parsedVersionURL.RawQuery = ""
			parsedVersionURL.Path += "/tags/" + url.PathEscape(app.Version.Fixed)
			versionURL = parsedVersionURL.String()
		} else if parsedVersionURL.RawQuery == "" && strings.HasSuffix(parsedVersionURL.Path, "/releases") {
			parsedVersionURL.Path += "/latest"
			versionURL = parsedVersionURL.String()
		}

		request, requestErr := http.NewRequest(http.MethodGet, versionURL, nil)
		if requestErr != nil {
			fmt.Println(aurora.Red("Failed:"), appID, "version request error:", requestErr)
			failCount++
			continue
		}

		if githubToken != "" && strings.HasSuffix(request.URL.Hostname(), "github.com") {
			request.Header.Set("Authorization", "Bearer "+githubToken)
		}

		response, httpErr := http.DefaultClient.Do(request)
		if httpErr != nil {
			fmt.Println(aurora.Red("Failed:"), appID, "version fetch error:", httpErr)
			failCount++
			continue
		}

		if response.StatusCode < 200 || response.StatusCode >= 300 {
			fmt.Println(aurora.Red("Failed:"), appID, "version fetch status:", response.Status)
			response.Body.Close()
			failCount++
			continue
		}

		body, readErr := io.ReadAll(response.Body)
		response.Body.Close()
		if readErr != nil {
			fmt.Println(aurora.Red("Failed:"), appID, "version body read error:", readErr)
			failCount++
			continue
		}

		if version == "" {
			matches := versionRegex.FindStringSubmatch(string(body))
			if len(matches) == 0 {
				fmt.Println(aurora.Red("Failed:"), appID, "version regex found no match")
				failCount++
				continue
			}

			version = matches[0]
			if len(matches) > 1 {
				version = matches[1]
			}
		}

		release := githubRelease{}
		if unmarshalErr := json.Unmarshal(body, &release); unmarshalErr != nil {
			releases := []githubRelease{}
			if unmarshalArrayErr := json.Unmarshal(body, &releases); unmarshalArrayErr != nil {
				fmt.Println(aurora.Red("Failed:"), appID, "invalid release payload:", unmarshalErr)
				failCount++
				continue
			}

			if len(releases) == 0 {
				fmt.Println(aurora.Red("Failed:"), appID, "release payload has no entries")
				failCount++
				continue
			}

			releaseFound := false
			for _, current := range releases {
				if current.TagName == version {
					release = current
					releaseFound = true
					break
				}

				if versionRegex == nil {
					continue
				}

				tagPayload := fmt.Sprintf("\"tag_name\":\"%s\"", current.TagName)
				tagMatches := versionRegex.FindStringSubmatch(tagPayload)
				if len(tagMatches) == 0 {
					tagPayload = fmt.Sprintf("\"tag_name\": \"%s\"", current.TagName)
					tagMatches = versionRegex.FindStringSubmatch(tagPayload)
				}

				if len(tagMatches) == 0 {
					continue
				}

				candidateVersion := tagMatches[0]
				if len(tagMatches) > 1 {
					candidateVersion = tagMatches[1]
				}

				if candidateVersion == version {
					release = current
					releaseFound = true
					break
				}
			}

			if !releaseFound {
				fmt.Println(aurora.Red("Failed:"), appID, "no release matched resolved version")
				failCount++
				continue
			}
		}

		architecture := map[string]map[string]string{}
		archFailed := false
		for arch, archConfig := range app.Autoupdate.Architecture {
			resolvedURL := strings.ReplaceAll(archConfig.URL, "$version", version)
			assetURL := strings.SplitN(resolvedURL, "#", 2)[0]

			asset := githubReleaseAsset{}
			assetFound := false
			for _, current := range release.Assets {
				if current.BrowserDownloadURL == assetURL {
					asset = current
					assetFound = true
					break
				}
			}

			if !assetFound {
				fmt.Println(aurora.Red("Failed:"), appID, "no matching release asset for", arch)
				failCount++
				archFailed = true
				break
			}

			hashValue := ""
			if strings.HasPrefix(strings.ToLower(asset.Digest), "sha256:") {
				hashValue = strings.ToLower(strings.TrimSpace(strings.SplitN(asset.Digest, ":", 2)[1]))
			} else {
				computedHash, hashErr := getURLSHA256(assetURL, githubToken)
				if hashErr != nil {
					fmt.Println(aurora.Red("Failed:"), appID, "hash fetch error:", hashErr)
					failCount++
					archFailed = true
					break
				}

				hashValue = computedHash
			}

			architecture[arch] = map[string]string{
				"url":  resolvedURL,
				"hash": hashValue,
			}
		}

		if archFailed {
			continue
		}

		manifest := map[string]any{
			"version":      version,
			"description":  app.Description,
			"homepage":     app.Homepage,
			"bin":          app.Bin,
			"shortcuts":    app.Shortcuts,
			"extract_dir":  app.Extract,
			"pre_install":  app.Install,
			"architecture": architecture,
			"autoupdate":   app.Autoupdate,
		}

		manifestRaw, marshalErr := json.MarshalIndent(manifest, "", "  ")
		if marshalErr != nil {
			fmt.Println(aurora.Red("Failed:"), appID, "manifest encode error:", marshalErr)
			failCount++
			continue
		}

		outputPath := resolvedScoopDir + "/" + appID + ".json"
		writeErr := os.WriteFile(outputPath, append(manifestRaw, '\n'), 0o644)
		if writeErr != nil {
			fmt.Println(aurora.Red("Failed:"), appID, "manifest write error:", writeErr)
			failCount++
			continue
		}

		fmt.Println(aurora.Green("Updated:"), outputPath)
		successCount++
	}

	fmt.Println()
	fmt.Println(aurora.Green("Updated Scoop manifests:"), successCount)
	if failCount > 0 {
		fmt.Println(aurora.Red("Failed Scoop manifests:"), failCount)
		os.Exit(1)
	}
}

func getURLSHA256(rawURL string, githubToken string) (string, error) {
	request, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return "", err
	}

	if githubToken != "" && strings.HasSuffix(request.URL.Hostname(), "github.com") {
		request.Header.Set("Authorization", "Bearer "+githubToken)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("hash download status: %s", resp.Status)
	}

	var h hash.Hash = sha256.New()
	if _, err := io.Copy(h, resp.Body); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
