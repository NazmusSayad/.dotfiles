<img width="2855" alt="image" src="https://github.com/user-attachments/assets/f217bd38-32a5-4bee-a259-96d3f5fb6837" />

# Development Setup for Windows

This repository contains automation and configuration for provisioning a Windows developer workstation. It includes shell configs, PowerShell/Batch helpers, AutoHotkey tooling, and small Go utilities.

## Features & Capabilities

- **Windows Configuration**
  PowerShell scripts to apply system settings and remove/disable unwanted defaults.

- **Apps, Packages, and Runtimes Management**
  Winget install/upgrade helpers driven by `config/winget-apps.jsonc`.

- **Shell Experience**
  Configured Bash/Fish/Zsh + Starship, Windows Terminal settings, and common aliases.

- **Code Editor Configuration**
  Editor configs (e.g. Zed), plus helpers (e.g. code snippet cleanup).

- **Communication Optimization**
  Slack helpers (startup + status).

- **System Performance**
  Debloat/tuning scripts under `src/ps1/` (review before running).

## Getting Started

### Prerequisites

- Windows 10 or Windows 11
- Git (installed to clone the repository)
- Go (to compile the utilities; see `go.mod` for the version)
- MSYS2 (optional, for `pacman`-managed shells/tools)

### Installation Guide

1.  **Clone the Repository:**

    ```shell
    git clone https://github.com/NazmusSayad/.dotfiles.git
    ```

2.  **Install Dotfiles (symlink + PATH):**
    Run `__install-dotfiles.cmd` as Administrator. This links the repo to `%USERPROFILE%\.dotfiles` and adds `%USERPROFILE%\.dotfiles\.build\bin` to the system PATH.

3.  **Compile utilities:**
    Run `__compile.cmd`. This compiles:

    - Go utilities from `src/scripts/*` and `src/functions/*` into `.build/bin/*.exe`
    - AutoHotkey scripts via `src/compile-ahk/` into `.build/ahk/`

4.  **Optional setup scripts:**

    - `__install-config.cmd`: Git + pnpm config.
    - `__git-gpg.cmd`: Generate/configure a GPG key for Git signing (prints the armored public key).
    - `__install-msys2.cmd`: Install shells/tools via MSYS2 `pacman`.
    - `__install-start-menu.cmd`: Install start-menu entries (via `go run ./src/install-start-menu/main.go`).
    - `__windows-setup.cmd`: Runs every PowerShell script in `src/ps1/` (admin required; reboots at the end).

5.  **Use the tools:**
    After compilation and PATH setup, the compiled tools are available as `*.exe` in `.build/bin/` (folder name becomes the exe name), e.g.:

    - `winget-install.exe`, `winget-upgrade.exe`
    - `symlink-setup.exe`
    - `slack-status.exe`, `slack-startup.exe`
    - Git helpers like `gclean.exe`, `greset.exe`, `gp.exe`, etc.

## Repository Structure

- `.build/`: Build output (binaries and compiled AHK).
- `config/`: Configuration files for shells, standard apps, and `winget` packages.
- `src/`: Go utilities and PowerShell scripts.
  - `src/functions/`: Small command-like Go utilities (compiled to `.build/bin/*.exe`).
  - `src/scripts/`: Higher-level Go scripts (compiled to `.build/bin/*.exe`).
  - `src/ps1/`: Windows debloating and configuration scripts.
  - `src/compile-go/`: Compiles Go utilities into `.build/bin/`.
  - `src/compile-ahk/`: Compiles bundled AHK scripts into `.build/ahk/`.
- `__*`: Installation and utility scripts.

## ⚠️ Disclaimer

This repository contains scripts that modify system settings and remove default applications. Review all scripts (especially those in `src/ps1/`) before running them to ensure they align with your requirements.
