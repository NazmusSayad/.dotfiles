//go:build !windows

package helpers

import (
	"errors"
)

type Scope string

const (
	ScopeUser    Scope = "User"
	ScopeMachine Scope = "Machine"
)

var errNotSupported = errors.New("environment helpers are only supported on Windows")

func ReadEnv(scope Scope, name string) (string, error) {
	return "", errNotSupported
}

func WriteEnv(scope Scope, name, value string) (string, error) {
	return "", errNotSupported
}

func AddToEnvPath(scope Scope, paths ...string) (string, error) {
	return "", errNotSupported
}
