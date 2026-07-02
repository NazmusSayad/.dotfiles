![Preview](https://github.com/user-attachments/assets/424e9395-a2b4-4fca-803b-4339df9616fc)

# Development Setup (Windows & macOS)

A complete automation system for setting up and managing your developer workstation on Windows and macOS. This repository provides everything you need to configure your development environment, manage applications, automate daily tasks, and optimize your workflow.

## Recommendation

I don't want everyone to copy this entire repo and use it because it is my own opinionated setup, not something for everyone.

I recommend checking this repo, taking ideas from it, and implementing them in a way that fulfills your own requirements.

## Features & Capabilities

- **Windows System Configuration**
  Automatically configure Windows settings, remove bloatware, disable unnecessary services, and optimize your system for development work.

- **macOS System Configuration**
  Install packages via Homebrew using the included Brewfile for a consistent macOS development environment.

- **Application Management**
  Install and update development tools, applications, and packages using Winget/Scoop (Windows) or Homebrew (macOS) with simple configuration files.

- **Enhanced Shell Experience**
  Pre-configured shell environments (Bash, Fish, Zsh) with Starship prompt, Windows Terminal / Ghostty settings, and convenient command aliases for faster workflow.

- **Git Workflow Tools**
  Streamlined Git commands for common tasks like cloning repositories, checking out branches, pulling changes, and managing your repositories more efficiently.

- **Smart Slack Integration**
  Automatically start or stop Slack based on your work schedule. Set your office hours, weekends, and off days, and Slack will manage itself accordingly.

- **Automated Startup Tasks**
  Configure applications and scripts to run automatically on system startup, with support for both user and administrator privileges.

- **Code Editor Setup**
  Pre-configured VS Code settings, keybindings, extensions, snippets, and state sync helpers.

## Agent Skills

Add the `npm` agent skill with:

```shell
skills add NazmusSayad/.dotfiles/config/ai/skills -s npm
```

## Getting Started

### Prerequisites

Before you begin, make sure you have:

- **Windows 10/11** or **macOS** installed
- **Git** installed (to clone the repository)
- **Go** installed (version specified in `go.mod`)
- **MSYS2** (Windows, optional for Bash/Fish) or **Homebrew** (macOS)
- **Scoop** (Windows package tools)

### Installation Guide

Follow these steps to set up your development environment:

1.  **Clone the Repository**

    Open a terminal and run:

    ```shell
    git clone https://github.com/NazmusSayad/.dotfiles.git
    cd .dotfiles
    ```

2.  **Initial Setup**

    - **Windows**: Right-click `__install-dotfiles.cmd` and select "Run as Administrator".
    - **macOS**: Run `__setup-macos.sh` after reviewing it, or bootstrap via Homebrew + config (see `config/Brewfile` and shell configs).

    This will set up the dotfiles directory and add tools to your PATH.

3.  **Build All Utilities**

    Run `__compile.sh`. On Windows, run it from Bash/MSYS2. This creates the executable tools for daily use.

4.  **Configure Your Environment**

    Run `__install-config.cmd` (Windows) or `__install-config.sh` (macOS/Linux) to set up:
    - Git configuration (name, email, default settings)
    - Symbolic links for configuration files
    - Scheduled tasks / launch agents for automatic startup
    - Shells and development tools
    - Go environment variables

5.  **Optional: Platform System Configuration**

    ⚠️ **Important:** Review the scripts before running this step!

    - **Windows**: Run `__setup-windows.cmd` as Administrator to apply system settings, remove bloatware, disable services, and optimize performance (restarts automatically).
    - **macOS**: Run `__setup-macos.sh` after reviewing it, and use Homebrew to provision packages from `config/Brewfile`.

6.  **Optional: Additional Setup**
    - Windows: `__compile-ahk.cmd`, `__install-code.cmd`
    - macOS: review `config/Brewfile`, `__setup-macos.sh`, and shell configs

### Using the Tools

Once installed, you can use these commands from anywhere in your terminal:

**Package Management:**

- `packages-sync` / `psy` - Sync configured packages
- `winget-install` / `wgi` - Install configured Winget applications (Windows)
- `winget-upgrade` / `wgu` - Update configured Winget applications (Windows)
- `scoop-install` - Install configured Scoop packages (Windows)

**Git Helpers:**

- `git-clone` / `gc` - Clone repositories (supports GitHub shorthand)
- `git-pull` / `gp` - Quick git pull
- `git-pull-all` / `gpa` - Pull all repositories in a workspace
- `gpr` - Pull changes with rebase
- `gpm` - Pull changes with merge
- `github-pr-create` / `ghp` - Create GitHub pull requests

**Slack Management:**

- `slack-status` / `ss` - Change Slack auto-start behavior (Always/Work Hours/Disabled)
- Slack will automatically start/stop based on your configured work schedule

**System Setup:**

- `symlink-init` - Recreate all configuration file symlinks
- `msys-init` - Set up MSYS2 development environment (Windows)
- `windows-startup` - Run configured Windows startup tasks

## Customization

Most user settings live under `config/`. Update the relevant config files, then re-run the platform install-config step to refresh links, scheduled tasks/launch agents, shortcuts, and environment settings.

Common files to edit:

- `config/apps.yaml` - Applications and packages to install (Windows Winget/Scoop/MSYS2)
- `config/Brewfile` - Packages to install (macOS Homebrew)
- `config/symlink.jsonc` - Config files linked into your system
- `config/slack-status.jsonc` - Slack startup schedule
- `config/shell/` - Shell aliases, prompt, and terminal settings
- `config/vscode/` - VS Code settings, keybindings, extensions, snippets, and synced state

## ⚠️ Important Notes

- **Review Before Running:** The Windows setup scripts (`src/ps1-windows/`) will modify system settings and remove default Windows applications. Please review these scripts before running `__setup-windows.cmd` to ensure they match your preferences.

- **Administrator Rights:** Some Windows scripts require administrator privileges. macOS scripts may prompt for your password via sudo.

- **Backup First:** Consider backing up important data before running system modification scripts.
