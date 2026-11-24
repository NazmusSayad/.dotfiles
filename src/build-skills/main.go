package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"dotfiles/src/helpers"
	"dotfiles/src/utils"
)

func main() {
	bundleSkills(resetDir("@/.local/skills-bundle"))
	writeFleetSkill(resetDir("@/.local/skills"))
}

func resetDir(path string) string {
	dir := helpers.ResolvePath(path)
	if err := os.RemoveAll(dir); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(err)
	}
	return dir
}

func bundleSkills(outputDir string) {
	skillsDir := helpers.ResolvePath("@/agents/skills")
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if !entry.IsDir() || !utils.IsFileExists(filepath.Join(skillsDir, entry.Name(), "SKILL.md")) {
			continue
		}
		zipSkill(skillsDir, entry.Name(), filepath.Join(outputDir, entry.Name()+".zip"))
		fmt.Println("Bundled " + entry.Name() + " -> .local/skills-bundle/" + entry.Name() + ".zip")
	}
}

func zipSkill(skillsDir string, name string, zipPath string) {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		panic(err)
	}
	defer zipFile.Close()

	writer := zip.NewWriter(zipFile)
	err = filepath.WalkDir(filepath.Join(skillsDir, name), func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		rel, err := filepath.Rel(skillsDir, path)
		if err != nil {
			return err
		}
		zipEntry, err := writer.Create(filepath.ToSlash(rel))
		if err != nil {
			return err
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(zipEntry, file)
		return err
	})
	if err != nil {
		panic(err)
	}
	if err := writer.Close(); err != nil {
		panic(err)
	}
}

func writeFleetSkill(outputDir string) {
	machines := helpers.ReadMachines()
	current := helpers.CurrentMachine(machines)

	var skill strings.Builder
	skill.WriteString("---\n")
	skill.WriteString("name: fleet\n")
	skill.WriteString("description: Use when a task needs to reach or run something on another of the user's machines.\n")
	skill.WriteString("---\n\n")
	fmt.Fprintf(&skill, "You are on %s (%s). %s Run local commands directly.\n", current.ID, current.Name, current.Description)

	for _, machine := range machines {
		if machine.ID == current.ID {
			continue
		}
		fmt.Fprintf(&skill, "\n## %s (%s)\n\n%s\n\n", machine.ID, machine.Name, machine.Description)
		fmt.Fprintf(&skill, "- SSH: %s\n", sshCommand(machine))
		for _, note := range machine.Notes {
			fmt.Fprintf(&skill, "- %s\n", note)
		}
	}

	skill.WriteString("\n> Note: SSH goes through Cloudflare Access. If `ssh` prints a Cloudflare login URL and waits, ask the user to log in.\n")

	skillDir := filepath.Join(outputDir, "fleet")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skill.String()), 0o644); err != nil {
		panic(err)
	}
	fmt.Println("Wrote fleet for " + current.ID + " -> .local/skills/fleet/SKILL.md")
}

func sshCommand(machine helpers.Machine) string {
	switch machine.OS {
	case "windows":
		return "`ssh " + machine.SSHHost() + "`"
	case "linux", "darwin":
		return "`ssh " + machine.SSHHost() + ` 'bash -lc "<command>"'` + "`"
	default:
		panic("fleet: unsupported OS " + machine.OS + " for " + machine.ID)
	}
}
