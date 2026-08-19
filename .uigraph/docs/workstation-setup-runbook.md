# Workstation Setup Runbook

## Purpose

This repository provisions an opinionated developer workstation for Windows and macOS. It installs or configures tools, creates links from repository-managed files into user configuration directories, and builds Go utilities used by the setup.

## Prerequisites

- Windows 10/11 or macOS.
- Git.
- Go 1.24.5, as declared in `go.mod`.
- Homebrew on macOS, and Scoop/MSYS2 as applicable on Windows.
- Administrator or `sudo` access for operations that change system settings.

## Main entry points

- Windows initial setup: review and run `__install-dotfiles.cmd` as administrator.
- Windows configuration: `__install-config.cmd` sources `etc/git-config.sh` and `etc/shell-config.sh`, invokes `src/scripts/symlink-init`, installs Windows tasks, and installs Start menu entries.
- macOS configuration: review and run `__setup-macos.sh`. It sources the Git and shell configuration scripts, changes the login shell, applies `launchctl`, `spctl`, `mdutil`, battery, Finder, Dock, keyboard, locale, and Safari defaults, opens the Chrome profile, and restarts Dock, Finder, and SystemUIServer.
- Build utilities: run `__compile.sh` from Bash/MSYS2. It removes `./.build/bin` and runs `go run ./src/compile-scripts/main.go`.
- Windows AutoHotkey build: `__compile-ahk.cmd` is an additional platform-specific entry point documented by the README.

## Safety

Review platform scripts before execution. Windows PowerShell setup scripts modify system settings, remove applications or capabilities, and disable services. macOS setup uses `sudo` and changes system defaults. Back up important workstation configuration first.

## Verification

Run `go test ./...` from the repository root after code or setup changes. Verify generated links, scheduled tasks or launch agents, shells, and expected applications on the target platform.