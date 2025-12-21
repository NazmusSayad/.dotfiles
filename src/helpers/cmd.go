package helpers

import (
	"os"
	"os/exec"
)

type ExecCommandOptions struct {
	Command string
	Args    []string

	Exit     bool
	Silent   bool
	NoStdin  bool
	NoStdout bool
	NoStderr bool
}

func ExecNativeCommand(options ExecCommandOptions) error {
	if options.Silent {
		options.NoStdin = true
		options.NoStdout = true
		options.NoStderr = true
	}

	cmd := exec.Command(options.Command, options.Args...)

	if !options.NoStdin {
		cmd.Stdin = os.Stdin
	}

	if !options.NoStdout {
		cmd.Stdout = os.Stdout
	}

	if !options.NoStderr {
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()
	if err != nil && options.Exit {
		if ee, ok := err.(*exec.ExitError); ok {
			os.Exit(ee.ExitCode())
		} else {
			os.Exit(1)
		}
	}

	return err
}
