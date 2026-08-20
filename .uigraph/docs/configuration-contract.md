# Configuration contract

Configuration is read relative to the repository root through the `@/` path convention in Go helpers. JSON and JSONC files are parsed with `github.com/tidwall/jsonc`; YAML files are parsed with `gopkg.in/yaml.v3`.

## Package manifests

`config/apps.yaml` has three top-level lists:

- `msys2`: package names or objects with `id` and optional `repo`.
- `scoop`: package names, optionally prefixed as `<bucket>/<app>`, or objects with `id` and optional `bucket`. A bare package defaults to the `main` bucket.
- `winget`: package IDs or objects supporting `id`, `name`, `scope`, `version`, `installerType`, administrative-install and upgrade flags, interactive flags, and skip flags.

`config/Brewfile` is a Homebrew bundle manifest containing `brew` formulae and `cask` applications. It includes shells, CLI tools, browsers, Git, development tools, Docker/Colima, media tools, and fonts.

## Linking and schedules

`config/symlink.jsonc` is an array of mappings with `Source`, `Target`, and optional platform-specific `Target.Win` or `Target.Mac`; `Copy: true` requests copying instead of linking. Sources use `@/` for repository paths and may reference user directories.

`config/slack-status.jsonc` defines `OfficeTimeStart: 8`, `OfficeTimeFinish: 20`, Friday/Saturday weekends, and month-to-day off-day lists. The Slack startup utility consumes this schedule.

## Shell and editor configuration

`config/shell/` contains Bash, Fish, Starship, Ghostty, btop, direnv, and terminal configuration. `config/vscode/` contains settings, keybindings, extensions, snippets, and state. `config/ai/` contains configuration for Claude, Codex, OpenCode, Pi, and shared agent instructions.

## Environment and generated paths

`__install-dotfiles.cmd` sets user `DOTFILES_DIR` to the repository directory and adds `$DOTFILES_DIR\\.build\\bin` to the user PATH. Compiled utilities are expected in `.build/bin`. `etc/git-config.sh` writes global Git identity and behavior, including `core.editor`, `core.excludesfile`, line endings, pull strategy, and default branch.
