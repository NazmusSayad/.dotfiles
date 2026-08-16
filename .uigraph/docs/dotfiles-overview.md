# Dotfiles Overview

The `.dotfiles` repository ("Dotfiles") is a complete automation system for
setting up and managing a developer workstation on **Windows** and **macOS**.
It is an opinionated per-machine setup owned by the repository author and is
not intended for wholesale adoption; it is provided so others can take ideas
from it and adapt them to their own environment.

## Purpose

The project provides everything needed to configure a development environment,
manage applications, automate daily tasks, and optimize workflow. It is
implemented primarily in **Go** (module `dotfiles`, `go 1.24.5`) and ships a
set of command-line utilities that are compiled and installed on the
workstation.

## Capabilities

- **Windows system configuration** — configure Windows settings, remove
  bloatware, disable unnecessary services, and optimize for development work.
- **macOS system configuration** — install packages via Homebrew using the
  bundled `config/Brewfile`.
- **Application management** — install and update development tools,
  applications, and packages via Winget/Scoop (Windows) or Homebrew (macOS),
  driven by `config/apps.yaml`.
- **Enhanced shell experience** — pre-configured Bash and Fish shells with the
  Starship prompt, Windows Terminal / Ghostty settings, and command aliases.
- **Git workflow tools** — streamlined commands for cloning repos, checking
  out branches, pulling changes, and managing multiple repositories.
- **Smart Slack integration** — automatically start or stop Slack based on a
  configured work schedule (`config/slack-status.jsonc`).
- **Automated startup tasks** — configure applications and scripts to run on
  system startup, with user and administrator privilege support.
- **Code editor setup** — pre-configured VS Code settings, keybindings,
  extensions, snippets, and state-sync helpers (including syncing VS Code /
  Cursor UI state into their internal SQLite `ItemTable`).

## Repository layout

- `config/` — the source of truth configuration (apps, Brewfile, symlinks,
  shell, VS Code, AI tooling, OBS, and more).
- `src/` — Go source for helpers, utilities, and the command-line scripts.
- `src/scripts/` — individual CLI tools (one `main.go` per command).
- `src/helpers/` — shared libraries used by the scripts.
- Root scripts and `.cmd` files — bootstrapping, compilation, and install
  helpers (`__setup-macos.sh`, `__install-dotfiles.cmd`, `__compile.sh`).

## Getting started

1. Clone the repository.
2. Run the platform installer (`__install-dotfiles.cmd` on Windows,
   `__setup-macos.sh` on macOS) to set up the dotfiles directory and add tools
   to `PATH`.
3. Build all utilities with `__compile.sh`.
4. Configure the environment (git config, symlinks for config files, scheduled
   tasks / launch agents, shells, and Go environment).
5. Optionally apply platform system configuration scripts, reviewing them first.

## Important notes

- Windows setup scripts (`src/ps1-windows/`) modify system settings and remove
  default Windows applications; review before running.
- Some Windows scripts require administrator privileges; macOS scripts may
  prompt for a password via `sudo`.
- Back up important data before running system-modification scripts.