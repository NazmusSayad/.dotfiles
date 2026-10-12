package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"dotfiles/src/utils"
)

type ExecCommandOptions struct {
	Dir      string
	Env      []string
	ExtraEnv map[string]string

	Exit     bool
	Silent   bool
	NoStdin  bool
	NoStdout bool
	NoStderr bool

	NoWait       bool
	AsAdmin      bool
	AsGsudoUser  bool
	AsGsudoAdmin bool

	Simulate bool
	Detached bool
}

func ExecNativeCommand(args []string, options ...ExecCommandOptions) error {
	if len(options) > 1 {
		panic("only one options struct is allowed")
	}

	opts := ExecCommandOptions{}
	if len(options) == 1 {
		opts = options[0]
	}

	if len(args) == 0 || strings.TrimSpace(args[0]) == "" {
		return fmt.Errorf("command is required")
	}

	if opts.AsGsudoUser && opts.AsAdmin {
		return fmt.Errorf("cannot run as gsudo user and admin at the same time")
	}

	commandArgs := append([]string(nil), args...)

	// On Windows, Simulate means executing through cmd.exe so commands
	// that depend on shell behavior are handled consistently.
	if opts.Simulate && runtime.GOOS == "windows" {
		commandLine := buildWindowsCommandLine(commandArgs)
		commandArgs = []string{"cmd", "/d", "/s", "/c", commandLine}
	}

	isAlreadyAsAdmin := isRunningAsAdmin()

	if opts.AsAdmin && !isAlreadyAsAdmin {
		commandArgs, err := prependElevatedCommand(commandArgs, opts)
		if err != nil {
			if opts.Exit {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			return err
		}
	}

	if opts.AsGsudoUser && isAlreadyAsAdmin {
		if !utils.IsCommandInPath("gsudo") {
			err := fmt.Errorf("gsudo not found in PATH")
			if opts.Exit {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			return err
		}

		commandArgs = append([]string{"gsudo", "--integrity", "Medium"}, commandArgs...)
	}

	cmd := exec.Command(commandArgs[0], commandArgs[1:]...)
	applySysProcAttr(cmd, opts.Detached)

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

	if opts.Dir != "" {
		cmd.Dir = opts.Dir
	}

	cmd.Env = buildEnvironment(opts.Env, opts.ExtraEnv)

	var err error

	if opts.NoWait {
		err = cmd.Start()
	} else {
		err = cmd.Run()
	}

	if err != nil && opts.Exit {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}

		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return err
}

func prependElevatedCommand(args []string, opts ExecCommandOptions) ([]string, error) {
	if opts.AsGsudoAdmin {
		if !utils.IsCommandInPath("gsudo") {
			return nil, fmt.Errorf("gsudo not found in PATH")
		}

		return append([]string{"gsudo"}, args...), nil
	}

	if !utils.IsCommandInPath("sudo") {
		return nil, fmt.Errorf("sudo not found in PATH")
	}

	return append([]string{"sudo"}, args...), nil
}

func buildEnvironment(baseEnv []string, extraEnv map[string]string) []string {
	var env []string

	if len(baseEnv) > 0 {
		env = append(env, baseEnv...)
	} else {
		env = append(env, os.Environ()...)
	}

	if len(extraEnv) == 0 {
		return env
	}

	// Environment variables are case-insensitive on Windows.
	caseInsensitive := runtime.GOOS == "windows"

	for key, value := range extraEnv {
		if strings.TrimSpace(key) == "" {
			continue
		}

		prefix := key + "="
		replaced := false

		for i, entry := range env {
			entryKey, _, ok := strings.Cut(entry, "=")
			if !ok {
				continue
			}

			keysEqual := entryKey == key
			if caseInsensitive {
				keysEqual = strings.EqualFold(entryKey, key)
			}

			if keysEqual {
				env[i] = prefix + value
				replaced = true
				break
			}
		}

		if !replaced {
			env = append(env, prefix+value)
		}
	}

	return env
}

// buildWindowsCommandLine prepares a command line for:
//
//	cmd.exe /d /s /c "<command line>"
//
// This is intentionally conservative because Simulate is only used internally
// by this project for commands that need Windows shell semantics.
func buildWindowsCommandLine(args []string) string {
	if len(args) == 0 {
		return ""
	}

	escaped := make([]string, len(args))

	for i, arg := range args {
		escaped[i] = quoteWindowsCommandArg(arg)
	}

	// /s /c requires the whole command to be surrounded by quotes when
	// special shell characters or spaces may be present.
	return `"` + strings.Join(escaped, " ") + `"`
}

func quoteWindowsCommandArg(arg string) string {
	if arg == "" {
		return `""`
	}

	needsQuotes := strings.ContainsAny(arg, " \t\"&()[]{}^=;!'+,`~|<>")
	if !needsQuotes {
		return arg
	}

	var b strings.Builder
	b.WriteByte('"')

	backslashes := 0

	for _, r := range arg {
		switch r {
		case '\\':
			backslashes++
			b.WriteRune(r)

		case '"':
			// Escape all preceding backslashes and the quote itself.
			b.WriteString(strings.Repeat(`\`, backslashes))
			b.WriteString(`\"`)
			backslashes = 0

		default:
			backslashes = 0
			b.WriteRune(r)
		}
	}

	// Escape trailing backslashes before the closing quote.
	if backslashes > 0 {
		b.WriteString(strings.Repeat(`\`, backslashes))
	}

	b.WriteByte('"')

	return b.String()
}
