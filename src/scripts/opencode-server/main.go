package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
)

type opencodeServerConfig struct {
	Port   int    `json:"port"`
	Tunnel string `json:"tunnel"`
}

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "Usage: opencode-server [stop]")
		os.Exit(1)
	}

	if len(os.Args) == 2 {
		if os.Args[1] != "stop" {
			fmt.Fprintln(os.Stderr, "Usage: opencode-server [stop]")
			os.Exit(1)
		}

		stopServer()
		return
	}

	startServer()
}

func startServer() {
	config := helpers.ReadConfig[opencodeServerConfig]("@/config/ai/opencode-server.jsonc")

	executablePath, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to resolve executable path:", err)
		os.Exit(1)
	}

	executablePath, err = filepath.Abs(executablePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to resolve executable path:", err)
		os.Exit(1)
	}

	runningPID, err := helpers.FindOtherPIDByExecutablePath(executablePath, int32(os.Getpid()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to inspect running processes:", err)
		os.Exit(1)
	}

	if runningPID != 0 {
		fmt.Println(aurora.Faint("> OpenCode server is already running"))
		return
	}

	opencodeCommand := exec.Command("opencode", "serve", "--port", fmt.Sprint(config.Port))
	opencodeCommand.Stdout = os.Stdout
	opencodeCommand.Stderr = os.Stderr
	if err := opencodeCommand.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to start opencode:", err)
		os.Exit(1)
	}

	cloudflaredCommand := exec.Command("cloudflared", "tunnel", "run", config.Tunnel)
	cloudflaredCommand.Stdout = os.Stdout
	cloudflaredCommand.Stderr = os.Stderr
	if err := cloudflaredCommand.Start(); err != nil {
		_ = helpers.TerminateProcessTree(int32(opencodeCommand.Process.Pid))
		fmt.Fprintln(os.Stderr, "failed to start cloudflared:", err)
		os.Exit(1)
	}

	fmt.Println(aurora.Faint("> OpenCode server started"))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	_ = helpers.TerminateProcessTree(int32(opencodeCommand.Process.Pid))
	_ = helpers.TerminateProcessTree(int32(cloudflaredCommand.Process.Pid))
}

func stopServer() {
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to resolve executable path:", err)
		os.Exit(1)
	}

	executablePath, err = filepath.Abs(executablePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to resolve executable path:", err)
		os.Exit(1)
	}

	runningPID, err := helpers.FindOtherPIDByExecutablePath(executablePath, int32(os.Getpid()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to inspect running processes:", err)
		os.Exit(1)
	}

	if runningPID == 0 {
		fmt.Println(aurora.Faint("> OpenCode server is not running"))
		return
	}

	if err := helpers.TerminateProcessTree(runningPID); err != nil {
		fmt.Fprintln(os.Stderr, "failed to stop opencode-server:", err)
		os.Exit(1)
	}

	fmt.Println(aurora.Faint("> OpenCode server stopped"))
}
