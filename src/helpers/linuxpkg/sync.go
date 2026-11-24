package linuxpkg

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
)

func Sync() {
	config := readConfig()
	lock := readLock()

	fmt.Println("✘", aurora.Faint("Uninstalling apt and deb packages..."))
	removeAptApps(config, lock)
	fmt.Println()

	fmt.Println("✘", aurora.Faint("Uninstalling snap packages..."))
	removeSnapApps(config, lock)
	fmt.Println()

	fmt.Println("✘", aurora.Faint("Autoremove apt dependencies..."))
	runAsAdmin("apt-get", "autoremove", "--purge", "-y")
	fmt.Println()

	fmt.Println("△", aurora.Faint("Updating apt..."))
	runAsAdmin("apt-get", "update")
	fmt.Println()

	fmt.Println("▼", aurora.Faint("Installing apt packages..."))
	installAptApps(config, lock)
	fmt.Println()

	fmt.Println("▼", aurora.Faint("Installing snap packages..."))
	installSnapApps(config, lock)
	fmt.Println()

	fmt.Println("▼", aurora.Faint("Installing deb packages..."))
	installDebApps(config, lock)
	fmt.Println()

	fmt.Println("△", aurora.Faint("Upgrading apt packages..."))
	runAsAdmin("apt-get", "upgrade", "-y")
	fmt.Println()

	fmt.Println("△", aurora.Faint("Upgrading snap packages..."))
	runAsAdmin("snap", "refresh")
	fmt.Println()
}

func removeAptApps(config appsConfig, lock appsLock) {
	removedAptIDs := []string{}
	for id := range lock.Apt {
		if !slices.Contains(config.Apt, id) {
			removedAptIDs = append(removedAptIDs, id)
		}
	}

	removedDebIDs := []string{}
	for id := range lock.Deb {
		if !slices.ContainsFunc(config.Deb, func(app debApp) bool { return app.ID == id }) {
			removedDebIDs = append(removedDebIDs, id)
		}
	}

	installedIDs := []string{}
	for _, id := range slices.Concat(removedAptIDs, removedDebIDs) {
		if getDpkgVersion(id) != "" {
			installedIDs = append(installedIDs, id)
		}
	}

	// Marked auto so the autoremove step purges them, keeping any still needed by other packages
	if len(installedIDs) > 0 {
		fmt.Println(aurora.Faint("- Removing " + strings.Join(installedIDs, ", ")))
		if runAsAdmin(append([]string{"apt-mark", "auto"}, installedIDs...)...) != nil {
			return
		}
	}

	for _, id := range removedAptIDs {
		delete(lock.Apt, id)
	}

	for _, id := range removedDebIDs {
		for _, source := range lock.Deb[id].AptSources {
			fmt.Println(aurora.Faint("- Removing apt source " + source))
			runAsAdmin("rm", "-f", source)
		}
		delete(lock.Deb, id)
	}

	writeLock(lock)
}

func removeSnapApps(config appsConfig, lock appsLock) {
	for id := range lock.Snap {
		if slices.ContainsFunc(config.Snap, func(app snapApp) bool { return app.ID == id }) {
			continue
		}

		if isSnapInstalled(id) {
			fmt.Println(aurora.Faint("- Removing snap " + id))
			if runAsAdmin("snap", "remove", id) != nil {
				continue
			}
		}

		delete(lock.Snap, id)
	}

	writeLock(lock)
}

func installAptApps(config appsConfig, lock appsLock) {
	missingIDs := []string{}
	for _, id := range config.Apt {
		if getDpkgVersion(id) == "" {
			missingIDs = append(missingIDs, id)
		}
	}

	if runAsAdmin(append([]string{"apt-get", "install", "-y"}, config.Apt...)...) != nil {
		return
	}

	for _, id := range missingIDs {
		lock.Apt[id] = lockEntry{
			Version:     getDpkgVersion(id),
			InstalledAt: time.Now().Format(time.RFC3339),
		}
	}

	writeLock(lock)
}

func installSnapApps(config appsConfig, lock appsLock) {
	for _, app := range config.Snap {
		if isSnapInstalled(app.ID) {
			continue
		}

		args := []string{"snap", "install", app.ID}
		if app.Channel != "" {
			args = append(args, "--channel="+app.Channel)
		}
		if app.Classic {
			args = append(args, "--classic")
		}

		fmt.Println(aurora.Faint("- Installing snap " + app.ID))
		if runAsAdmin(args...) != nil {
			continue
		}

		lock.Snap[app.ID] = lockEntry{
			Channel:     app.Channel,
			Classic:     app.Classic,
			InstalledAt: time.Now().Format(time.RFC3339),
		}
		writeLock(lock)
	}
}

func installDebApps(config appsConfig, lock appsLock) {
	for _, app := range config.Deb {
		installedVersion := getDpkgVersion(app.ID)
		url := ""

		switch {
		case app.URL != "" && app.GitHub == "" && app.Asset == "":
			if installedVersion != "" {
				continue
			}
			url = app.URL
		case app.URL == "" && app.GitHub != "" && app.Asset != "":
			release, err := helpers.GetLatestGithubRelease(app.GitHub)
			if err != nil {
				fmt.Println(aurora.Red("Failed to get latest release of " + app.GitHub + ": " + err.Error()))
				continue
			}
			if installedVersion == strings.TrimPrefix(release.TagName, "v") {
				continue
			}

			assetPattern, err := regexp.Compile(app.Asset)
			if err != nil {
				fmt.Println(aurora.Red("Invalid asset pattern for " + app.ID + ": " + err.Error()))
				continue
			}
			assetIndex := slices.IndexFunc(release.Assets, func(asset helpers.Asset) bool {
				return assetPattern.MatchString(asset.Name)
			})
			if assetIndex == -1 {
				fmt.Println(aurora.Red("No asset of " + app.GitHub + " " + release.TagName + " matches " + app.Asset))
				continue
			}
			url = release.Assets[assetIndex].BrowserDownloadURL
		default:
			fmt.Println(aurora.Red("Invalid deb app " + app.ID + ": expected either url, or github with asset"))
			continue
		}

		fmt.Println(aurora.Faint("- Installing deb " + app.ID + " from " + url))
		sourcesBefore := listAptSources()
		if err := installDebFromURL(url); err != nil {
			fmt.Println(aurora.Red("Failed to install deb " + app.ID + ": " + err.Error()))
			continue
		}

		version := getDpkgVersion(app.ID)
		if version == "" {
			fmt.Println(aurora.Red("Installed deb from " + url + " does not provide package " + app.ID))
			continue
		}

		aptSources := lock.Deb[app.ID].AptSources
		for _, source := range listAptSources() {
			if !slices.Contains(sourcesBefore, source) && !slices.Contains(aptSources, source) {
				aptSources = append(aptSources, source)
			}
		}

		lock.Deb[app.ID] = lockEntry{
			Version:     version,
			URL:         url,
			GitHub:      app.GitHub,
			AptSources:  aptSources,
			InstalledAt: time.Now().Format(time.RFC3339),
		}
		writeLock(lock)
	}
}
