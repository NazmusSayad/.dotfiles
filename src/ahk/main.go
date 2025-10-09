package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var scriptPrefix = "___$"
var allowedScripts = map[string]struct{}{
	"AHK-Macro":         {},
	"AHK-WindowManager": {},
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	ahkScriptsDir := filepath.Join(cwd, "src", "ahk", "scripts")
	ahk2ExeBin := filepath.Join(cwd, "src", "ahk", "bin", "Ahk2Exe.exe")
	ahkCompilerBin := filepath.Join(cwd, "src", "ahk", "bin", "AutoHotkey64.exe")

	entries, err := os.ReadDir(ahkScriptsDir)
	if err != nil {
		panic(err)
	}

	compiledAhkScripts := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".ahk") {
			continue
		}
		fileName := strings.TrimSuffix(entry.Name(), ".ahk")
		inPath := filepath.Join(ahkScriptsDir, entry.Name())
		outPath := filepath.Join(cwd, scriptPrefix+fileName+".exe")
		iconPath := filepath.Join(ahkScriptsDir, fileName+".ico")
		spawnArgs := []string{"/base", ahkCompilerBin, "/in", inPath, "/out", outPath}
		if _, err := os.Stat(iconPath); err == nil {
			spawnArgs = append(spawnArgs, "/icon", iconPath)
		}
		cmd := exec.Command(ahk2ExeBin, spawnArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		_ = cmd.Run()
		fmt.Printf("Compiled: %s\n", entry.Name())
		compiledAhkScripts = append(compiledAhkScripts, outPath)
	}

	vbsEntries := make([]string, 0)
	vbsEntries = append(vbsEntries, formatProgramForVBS([]string{"cmd.exe", "/c"}))

	adminPrograms := make([]string, 0)
	for _, file := range compiledAhkScripts {
		baseName := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		baseNameId := strings.TrimPrefix(baseName, scriptPrefix)

		if _, ok := allowedScripts[baseNameId]; ok {
			adminPrograms = append(adminPrograms, file)
		}
	}
	adminPrograms = append(adminPrograms, `C:\Program Files\ShareX\ShareX.exe`)

	for _, exePath := range adminPrograms {
		vbsEntries = append(vbsEntries, formatProgramForVBSAsAdmin(exePath))
	}

	vbsEntries = append(vbsEntries, formatProgramForVBS([]string{"gpg", "--list-keys"}))

	launchVBS := filepath.Join(cwd, "___launch.vbs")
	if err := os.WriteFile(launchVBS, []byte(strings.Join(vbsEntries, "\n\n")), 0644); err != nil {
		panic(err)
	}

	fmt.Println("Compilation complete")

	taskXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.2" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <RegistrationInfo>
    <Date>%s</Date>
    <Author>DESKTOP-VP9L1TR\Sayad</Author>
    <URI>\#START</URI>
  </RegistrationInfo>
  <Triggers>
    <BootTrigger>
      <Enabled>true</Enabled>
    </BootTrigger>
    <LogonTrigger>
      <Enabled>true</Enabled>
    </LogonTrigger>
  </Triggers>
  <Principals>
    <Principal id="Author">
      <GroupId>S-1-5-32-544</GroupId>
      <RunLevel>HighestAvailable</RunLevel>
    </Principal>
  </Principals>
  <Settings>
    <MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>
    <DisallowStartIfOnBatteries>true</DisallowStartIfOnBatteries>
    <StopIfGoingOnBatteries>true</StopIfGoingOnBatteries>
    <AllowHardTerminate>true</AllowHardTerminate>
    <StartWhenAvailable>false</StartWhenAvailable>
    <RunOnlyIfNetworkAvailable>false</RunOnlyIfNetworkAvailable>
    <IdleSettings>
      <StopOnIdleEnd>true</StopOnIdleEnd>
      <RestartOnIdle>false</RestartOnIdle>
    </IdleSettings>
    <AllowStartOnDemand>true</AllowStartOnDemand>
    <Enabled>true</Enabled>
    <Hidden>true</Hidden>
    <RunOnlyIfIdle>false</RunOnlyIfIdle>
    <WakeToRun>false</WakeToRun>
    <ExecutionTimeLimit>PT72H</ExecutionTimeLimit>
    <Priority>7</Priority>
  </Settings>
  <Actions Context="Author">
    <Exec>
      <Command>%s</Command>
    </Exec>
  </Actions>
</Task>`, isoTimestamp(), launchVBS)

	if err := os.WriteFile(filepath.Join(cwd, "___task-init.xml"), []byte(taskXML), 0644); err != nil {
		panic(err)
	}
}

func formatProgramForVBS(program []string) string {
	if len(program) == 0 {
		return ""
	}
	exe := program[0]
	rest := program[1:]
	literal := fmt.Sprintf(`"""%s"""`, exe)
	if len(rest) > 0 {
		literal = fmt.Sprintf(`"""%s"" %s"`, exe, strings.Join(rest, ","))
	}
	lines := []string{
		"Set WshShell = CreateObject(\"WScript.Shell\")",
		fmt.Sprintf("WshShell.Run %s, 0, False", literal),
		"Set WshShell = Nothing",
	}
	return strings.Join(lines, "\n")
}

func formatProgramForVBSAsAdmin(program string) string {
	lines := []string{
		"Set UAC = CreateObject(\"Shell.Application\")",
		fmt.Sprintf("UAC.ShellExecute \"%s\", \"\", \"\", \"runas\", 0", program),
		"Set UAC = Nothing",
	}
	return strings.Join(lines, "\n")
}

func isoTimestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
}
