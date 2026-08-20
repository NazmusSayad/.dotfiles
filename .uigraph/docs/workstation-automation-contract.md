# Workstation Automation Contract

## Identity and scope

`Dotfiles` is a cross-platform workstation automation repository. It manages shell and Git configuration, application/package manifests, editor settings, symlinks, startup tasks, and platform-specific system settings. The repository is opinionated for the author's Windows and macOS environments; it is not a generic application service.

## Prerequisites

The README identifies Windows 10/11 or macOS, Git, Go 1.24.5 from `go.mod`, and platform tooling: MSYS2 or Homebrew, plus Scoop for Windows package workflows.

## Entry points and commands

- Windows one-time setup: run `__install-dotfiles.cmd` as Administrator.
- macOS setup: review and run `__setup-macos.sh`; it sources `etc/git-config.sh` and `etc/shell-config.sh`, changes macOS defaults, opens the Chrome profile, and restarts Dock/Finder/SystemUIServer.
- Windows system configuration: run `__setup-windows.cmd` as Administrator. It executes every `src/ps1-windows/*.ps1` script and restarts Windows.
- Configuration installation: `__install-config.cmd` runs Git/shell setup, `symlink-init`, Windows task installation, and start-menu installation.
- Build utilities: run `__compile.sh` from Bash/MSYS2. It removes `./.build/bin` and invokes `go run ./src/compile-scripts/main.go`.
- Optional utilities include `__compile-ahk.cmd` and `__install-code.cmd`.

The repository's `AGENTS.md` prohibits running `go build`, `go run`, repository sync/update scripts, or generated executables during automated work unless explicitly requested by the user.

## Checked-in configuration contract

- `config/apps.yaml` declares MSYS2, Scoop, and Winget packages, including interactive and administrative install flags.
- `config/Brewfile` declares Homebrew packages for macOS.
- `config/symlink.jsonc` maps repository sources to user/application destinations. It includes shell files, AI tools, terminals, OBS, VS Code, and other application state; some mappings are platform-specific and some use copy semantics.
- `config/slack-status.jsonc` defines office hours (08:00–20:00), weekend names, and month/day off-days used by Slack startup logic.
- `config/shell/`, `config/vscode/`, `config/ai/`, and application-specific directories contain the target configuration payloads.
- `.env` exists in the checkout and is treated as local input; no application source read during onboarding established a complete environment-variable schema. Secrets must not be copied into artifacts.

## Filesystem and operating-system boundaries

Setup and helpers read checked-in configuration and write symlinks/copies into paths under `$HOME`, `~/.config`, `~/.codex`, `~/.claude*`, `~/.local/share`, Windows `%APPDATA%`/`%LOCALAPPDATA%`, and macOS `~/Library/Application Support`. They also create `.build/bin` utilities. Windows scripts invoke PowerShell, Windows Task Scheduler, start-menu facilities, and a forced restart; macOS setup invokes `sudo`, `defaults`, `launchctl`, `mdutil`, `chflags`, `open`, and process termination commands.

## External integrations

- `src/helpers/github-release.go` calls `https://api.github.com/repos/{owner}/{repo}/releases/latest`, then downloads a matching asset URL over HTTPS. HTTP failures and non-200 responses are returned as errors.
- `src/helpers/opencode/fetch.go` calls `https://models.dev/api.json` and `https://openrouter.ai/api/v1/models`; OpenRouter may receive `Authorization: Bearer ...` when `OPENROUTER_API_KEY` resolves to a non-empty key. Responses are cached in memory for the process and errors are returned.
- Slack control is local process management: it discovers Slack through `scoop which slack`, checks the process with PowerShell, and starts/stops it with platform-native commands. It does not call a Slack network API.
- GitHub authentication is obtained by invoking the external `gh auth token` command; an unavailable command yields an empty token.

## Build and test evidence

The module is Go 1.24.5 and includes a test file at `src/helpers/json-patch_test.go`. No repository-wide test command, API contract, database schema, deployment manifest, or test-pack contract was found, so UiGraph test packs and database/API artifacts are intentionally not declared.
