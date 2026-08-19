# Go Utilities and Test Guide

## Module and build surface

The module is `dotfiles` and targets Go 1.24.5. `src/` contains reusable helpers, constants, utilities, Windows setup programs, compilation programs, and many standalone daily-use commands under `src/scripts/`.

`__compile.sh` runs `go run ./src/compile-scripts/main.go` after clearing `./.build/bin`. Windows setup entry points invoke Go programs with `go run`, including symlink initialization and task or Start menu installation.

## Runtime integrations

The code crosses these boundaries:

- Operating-system commands through `os/exec`, including PowerShell, `taskkill`, `pkill`, Scoop, native application commands, and other platform tools.
- Local files and user paths through filesystem helpers, configuration readers, symlink management, and generated build output.
- GitHub's releases API and release asset URLs in `src/helpers/github-release.go`.
- Configurable OpenAI-compatible model endpoints and `https://models.dev/api.json` in `src/helpers/opencode/fetch.go`; OpenRouter is a built-in provider and can use `OPENROUTER_API_KEY`.
- A locally installed Slack process through PowerShell/Scoop lookup and platform process commands.

Network calls are handled as direct HTTP requests. The external service is not required for ordinary local configuration unless the invoked utility uses that integration.

## Tests

The repository contains Go tests in `src/helpers/json-patch_test.go`, covering JSON merge preservation, additions, updates, removals, arrays, nested objects, and invalid input. Run:

```sh
go test ./...
```

The command compiles all Go packages and executes the checked-in tests. No API specification, database migration, or service-side test framework is present in the repository.