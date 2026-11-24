package linuxpkg

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
)

const aptSourcesDir = "/etc/apt/sources.list.d"

func runAsAdmin(args ...string) error {
	err := helpers.ExecNativeCommand(args, helpers.ExecCommandOptions{AsAdmin: true})
	if err != nil {
		fmt.Println(aurora.Red("Failed: " + strings.Join(args, " ") + ": " + err.Error()))
	}

	return err
}

func getDpkgVersion(id string) string {
	output, err := exec.Command("dpkg-query", "-W", "-f=${db:Status-Status} ${Version}\n", id).Output()
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		return ""
	}
	if err != nil {
		panic("dpkg-query failed for " + id + ": " + err.Error())
	}

	firstLine, _, _ := strings.Cut(string(output), "\n")
	status, version, _ := strings.Cut(firstLine, " ")
	if status != "installed" {
		return ""
	}

	return version
}

func isSnapInstalled(id string) bool {
	err := exec.Command("snap", "list", id).Run()
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		return false
	}
	if err != nil {
		panic("snap list failed for " + id + ": " + err.Error())
	}

	return true
}

func listAptSources() []string {
	entries, err := os.ReadDir(aptSourcesDir)
	if err != nil {
		panic("failed to read " + aptSourcesDir + ": " + err.Error())
	}

	sources := []string{}
	for _, entry := range entries {
		sources = append(sources, filepath.Join(aptSourcesDir, entry.Name()))
	}

	return sources
}

func installDebFromURL(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status code %d", resp.StatusCode)
	}

	file, err := os.CreateTemp("", "*.deb")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	_, copyErr := io.Copy(file, resp.Body)
	closeErr := file.Close()
	if copyErr != nil {
		return copyErr
	}
	if closeErr != nil {
		return closeErr
	}

	// apt reads local debs as the _apt user, and CreateTemp makes the file 0600
	if err := os.Chmod(file.Name(), 0o644); err != nil {
		return err
	}

	return helpers.ExecNativeCommand(
		[]string{"apt-get", "install", "-y", file.Name()},
		helpers.ExecCommandOptions{AsAdmin: true},
	)
}
