# Configuration Contract

Repository-managed workstation state is primarily under `config/`. Setup utilities read these files and link or copy them to platform-specific user locations.

## Package inputs

- `config/apps.yaml` declares packages for MSYS2, Scoop, and Winget. It includes shells and developer tools, system utilities, media tools, Slack, fonts, VS Code, SDKs, Docker Desktop, and other workstation applications.
- `config/Brewfile` is the macOS Homebrew package input referenced by the README.
- `config/mise-config.toml` configures mise.

## Link and copy mappings

`config/symlink.jsonc` is a JSON-with-comments list of mappings. Each entry has a repository `Source`, one or more `Target` paths, and optional platform-specific `Target.Win` or `Target.Mac` values. `Copy: true` makes the setup copy rather than link that item.

Examples implemented in the mapping include shell profiles to `~/.bash_profile`, `~/.bashrc`, and Fish config; AI agent configuration to Claude, Codex, and OpenCode locations; Starship, direnv, mise, btop, lsd, Ghostty, Windows Terminal, OBS, Raycast, VS Code, and other application configuration.

Sources beginning with `@/` resolve within this repository. Other sources such as `~/.agents/skills` and `~/.claude/*` refer to existing user-managed paths and are not repository files.

## Other configuration domains

- `config/shell/` contains Bash, Fish, Starship, terminal, and command aliases.
- `config/vscode/` contains settings, keybindings, extensions, snippets, and synced state.
- `config/ai/` contains agent instructions and configuration for Claude, OpenCode, Codex, and pi.
- `config/slack-status.jsonc` controls Slack startup scheduling; the Go Slack helpers start or stop the application according to that configuration.
- `config/symlink.jsonc` and the setup scripts are the authoritative source for destination paths; destination paths are user or operating-system state outside the repository.