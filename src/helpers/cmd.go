package helpers

import (
	"os"
	"os/exec"
)

type ExecCommandOptions struct {
	Command string
	Args    []string
	Env     []string

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

	if len(options.Env) > 0 {
		cmd.Env = options.Env
	} else {
		cmd.Env = os.Environ()
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

func SimulateCommandAlias(alias []string) error {
	aliasCommand := alias[0]
	aliasArguments := alias[1:]
	scriptArguments := os.Args[1:]

	return ExecNativeCommand(ExecCommandOptions{
		Command: aliasCommand,
		Args:    append(aliasArguments, scriptArguments...),
		Exit:    true,
	})
}
