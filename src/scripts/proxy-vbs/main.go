package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: proxy-vbs <user|admin> <executable>")
		return
	}

	mode := os.Args[1]
	executable := os.Args[2]

	if mode != "user" && mode != "admin" {
		fmt.Println("Mode must be 'user' or 'admin'")
		return
	}

	absPath, err := exec.LookPath(executable)
	if err != nil {
		fmt.Println("Executable not found in PATH:", executable)
		return
	}

	fmt.Printf("Resolved executable: %s (mode: %s)\n", absPath, mode)

	program := strings.ReplaceAll(absPath, `"`, `""`)

	var vbscript string
	if mode == "admin" {
		vbscript = formatProgramForVBSAsAdmin(program)
	} else {
		vbscript = formatProgramForVBS(program, os.Args[3:])
	}

	fmt.Println(vbscript)

	tempFile, err := os.CreateTemp("", "*.vbs")
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		return
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString(vbscript); err != nil {
		fmt.Println("Error writing to temp file:", err)
		return
	}
	tempFile.Close()

	cmd := exec.Command("cscript", tempFile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func formatProgramForVBS(program string, programArgs []string) string {
	literal := fmt.Sprintf(`"""%s"""`, program)
	if len(programArgs) > 0 {
		literal = fmt.Sprintf(`"""%s"" %s"`, program, strings.Join(programArgs, ","))
	}
	lines := []string{
		"Set WshShell = CreateObject(\"WScript.Shell\")",
		fmt.Sprintf("WshShell.Run %s, 0, False", literal),
		"Set WshShell = Nothing",
	}
	return strings.Join(lines, "\n")
}

// TODO: Add support for program arguments
func formatProgramForVBSAsAdmin(program string) string {
	lines := []string{
		"Set UAC = CreateObject(\"Shell.Application\")",
		fmt.Sprintf("UAC.ShellExecute \"%s\", \"\", \"\", \"runas\", 0", program),
		"Set UAC = Nothing",
	}
	return strings.Join(lines, "\n")
}
