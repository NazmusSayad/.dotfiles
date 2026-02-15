package helpers

import (
	"os"
	"os/exec"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

type ExecCommandOptions struct {
	Env []string

	Exit     bool
	Silent   bool
	NoStdin  bool
	NoStdout bool
	NoStderr bool

	NoWait  bool
	AsUser  bool
	AsAdmin bool

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

	if opts.AsUser && opts.AsAdmin {
		panic("cannot run as user and admin at the same time")
	}

	if opts.Simulate {
		args = []string{"cmd", "/c", strings.Join(args, " ")}
	}

	isAlreadyAsAdmin := isRunningAsAdmin()

	if opts.AsAdmin && !isAlreadyAsAdmin {
		args = append([]string{"sudo"}, args...)
	}

	if opts.AsUser && isAlreadyAsAdmin {
		args = append([]string{"gsudo", "--integrity", "Medium"}, args...)
	}

	cmd := exec.Command(args[0], args[1:]...)

	if opts.Detached {
		cmd.SysProcAttr = &windows.SysProcAttr{
			HideWindow:       true,
			NoInheritHandles: true,
			CreationFlags:    windows.DETACHED_PROCESS | windows.CREATE_NEW_PROCESS_GROUP,
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
