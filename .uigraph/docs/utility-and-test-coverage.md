# Utility architecture and test coverage

## Go utility layer

The Go module is `dotfiles`, targeting Go 1.24.5. Code under `src/helpers` provides configuration loading, application manifest parsing, filesystem and path operations, command execution, GitHub authentication and release downloads, symlink handling, Windows task support, and JSON patch/merge behavior. `src/scripts` contains the compiled command entrypoints for package synchronization, Git helpers, Slack startup/status, symlink management, Windows environment/startup, and related utilities.

The helper library depends on Bubble Tea/Lip Gloss/Aurora for terminal interfaces and output, JSONC/YAML parsers for configuration, `modernc.org/sqlite` for SQLite support in the Go dependency graph, and OS/process libraries for platform integration.

## External boundaries

The code reads repository configuration files and user/home-directory state, writes generated binaries and symlinks, invokes local OS services and package managers, and accesses GitHub over HTTPS for latest-release metadata and assets. GitHub authentication is optionally read from the local `gh auth token` command; a failed lookup returns an empty token.

## Existing tests

`src/helpers/json-patch_test.go` is a Go unit-test suite for `MergeJSONObject`. It covers unchanged multiline and compact passthrough, key addition/update/removal, array replacement, preservation of surrounding formatting and raw unchanged values, nested-object updates, repeated no-change cases, and invalid JSON/object inputs. The repository does not declare a separate API test framework or integration-test suite.

Run the repository's tests with `go test ./...` after reviewing repository guidance. Do not run build or setup commands automatically; `AGENTS.md` specifically prohibits automatic `go build`, `go run`, and direct execution of repository update scripts.
