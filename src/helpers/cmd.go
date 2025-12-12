package helpers

import (
	"os"
	"os/exec"
)

type ExecCommandOptions struct {
	Command string
	Args    []string

	Exit bool

	Stdin  bool
	Stdout bool
	Stderr bool
}

func ExecNativeCommand(options ExecCommandOptions) error {
	cmd := exec.Command(options.Command, options.Args...)

	if options.Stdin {
		cmd.Stdin = os.Stdin
	}

	if options.Stdout {
		cmd.Stdout = os.Stdout
	}

	if options.Stderr {
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
