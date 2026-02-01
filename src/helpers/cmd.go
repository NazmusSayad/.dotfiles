package helpers

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type ExecCommandOptions struct {
	Env []string

	Exit     bool
	Silent   bool
	NoStdin  bool
	NoStdout bool
	NoStderr bool

	Simulate bool
}

func ExecNativeCommand(args []string, options ...ExecCommandOptions) error {
	opts := ExecCommandOptions{}
	if len(options) > 0 {
		opts = options[0]
	}

	command := args[0]
	if len(args) == 0 || command == "" {
		panic("command is required")
	}

	cmd := exec.Command(command, args[1:]...)
	if opts.Simulate {
		cmd = exec.Command("cmd", "/c", strings.Join(args, " "))
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	if opts.Silent {
		opts.NoStdin = true
		opts.NoStdout = true
		opts.NoStderr = true
	}

	if !opts.NoStdin {
		cmd.Stdin = os.Stdin
	}

	if !opts.NoStdout {
		cmd.Stdout = os.Stdout
	}

	if !opts.NoStderr {
		cmd.Stderr = os.Stderr
	}

	if len(opts.Env) > 0 {
		cmd.Env = opts.Env
	} else {
		cmd.Env = os.Environ()
	}

	err := cmd.Run()
	if err != nil && opts.Exit {
		if ee, ok := err.(*exec.ExitError); ok {
			os.Exit(ee.ExitCode())
		} else {
			os.Exit(1)
		}
	}

	return err
}
