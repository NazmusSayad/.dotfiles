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

	NoWait   bool
	AsAdmin  bool
	Simulate bool
	Detached bool
}

func ExecNativeCommand(args []string, options ...ExecCommandOptions) error {
	opts := ExecCommandOptions{}
	if len(options) > 0 {
		opts = options[0]
	} else if len(options) > 1 {
		panic("only one options struct is allowed")
	}

	if len(args) == 0 || args[0] == "" {
		panic("command is required")
	}

	cmd := exec.Command(args[0], args[1:]...)
	if opts.Simulate && opts.AsAdmin {
		cmd = exec.Command("sudo", "cmd", "/c", strings.Join(args, " "))
	} else if opts.Simulate {
		cmd = exec.Command("cmd", "/c", strings.Join(args, " "))
	} else if opts.AsAdmin {
		cmd = exec.Command("sudo", args...)
	}

	if opts.Detached {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags:    syscall.CREATE_NEW_PROCESS_GROUP | 0x00000008,
			NoInheritHandles: true,
			HideWindow:       true,
		}
	} else {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
	}

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

	var err error
	if opts.NoWait {
		err = cmd.Start()
	} else {
		err = cmd.Run()
	}

	if err != nil && opts.Exit {
		if ee, ok := err.(*exec.ExitError); ok {
			os.Exit(ee.ExitCode())
		} else {
			os.Exit(1)
		}
	}

	return err
}
