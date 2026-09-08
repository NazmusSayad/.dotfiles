//go:build windows

package helpers

func inheritOwnership(target string, createdDirs []string) {}

func inheritDirOwnership(dirs []string) {}

func ApplyUserOwnership(paths ...string) {}

func ValidateInvokingUser() error { return nil }
