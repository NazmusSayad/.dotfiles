<img width="2855" alt="image" src="https://github.com/user-attachments/assets/f217bd38-32a5-4bee-a259-96d3f5fb6837" />

# Development Setup for Windows

A complete automation system for setting up and managing your Windows developer workstation. This repository provides everything you need to configure your development environment, manage applications, automate daily tasks, and optimize your workflow.

## Features & Capabilities

- **Windows System Configuration**
  Automatically configure Windows settings, remove bloatware, disable unnecessary services, and optimize your system for development work.

- **Application Management**
  Install and update all your development tools, applications, and packages using Winget with a simple configuration file. Keep everything up-to-date with one command.

- **Enhanced Shell Experience**
  Pre-configured shell environments (Bash, Fish, Zsh) with Starship prompt, Windows Terminal settings, and convenient command aliases for faster workflow.

- **Git Workflow Tools**
  Streamlined Git commands for common tasks like cloning repositories, checking out branches, pulling changes, and managing your repositories more efficiently.

- **Smart Slack Integration**
  Automatically start or stop Slack based on your work schedule. Set your office hours, weekends, and off days, and Slack will manage itself accordingly.

- **Automated Startup Tasks**
  Configure applications and scripts to run automatically on system startup, with support for both user and administrator privileges.

- **Code Editor Setup**
  Pre-configured settings for popular editors like Zed, along with helpful utilities for managing your development environment.

## Getting Started

### Prerequisites

Before you begin, make sure you have:

- **Windows 10 or Windows 11** installed
- **Git** installed (to clone the repository)
- **Go** installed (version specified in `go.mod`)
- **MSYS2** (optional, if you want to use Bash/Fish/Zsh shells)

### Installation Guide

Follow these steps to set up your development environment:

1.  **Clone the Repository**

    Open PowerShell or Command Prompt and run:

    ```shell
    git clone https://github.com/NazmusSayad/.dotfiles.git
    cd .dotfiles
    ```

2.  **Initial Setup**

    Right-click `__install-dotfiles.cmd` and select "Run as Administrator". This will:
    - Set up the dotfiles directory
    - Add the tools to your system PATH so you can use them anywhere

3.  **Build All Utilities**

    Run `__compile.cmd` to build all the tools and scripts. This creates executable files that you'll use for daily tasks.

4.  **Configure Your Environment**

    Run `__install-config.cmd` to set up:
    - Git configuration (name, email, default settings)
    - Symbolic links for configuration files
    - Windows scheduled tasks for automatic startup
    - Start menu shortcuts for quick access
    - MSYS2 shells and development tools
    - Go environment variables

5.  **Optional: Windows System Configuration**

    ⚠️ **Important:** Review the scripts before running this step!

    Run `__windows-setup.cmd` as Administrator to:
    - Apply Windows system settings
    - Remove default bloatware applications
    - Disable unnecessary services
    - Optimize system performance

    **Note:** This will restart your computer automatically after completion.

6.  **Optional: Additional Setup**

    - `__git-gpg.cmd`: Set up GPG key for Git commit signing
    - `__install-start-menu.cmd`: Add shortcuts to Windows Start Menu

### Using the Tools

Once installed, you can use these commands from anywhere in your terminal:

**Package Management:**
- `winget-install` - Install all configured applications
- `winget-upgrade` - Update all installed packages

**Git Helpers:**
- `c` - Clone repositories (supports GitHub shorthand)
- `gc` - Checkout branches (creates if doesn't exist)
- `gpr` - Pull changes with rebase
- `gpm` - Pull changes with merge
- `gp` - Quick git pull
- `gds` - Git diff with statistics

**Slack Management:**
- `slack-status` - Change Slack auto-start behavior (Always/Work Hours/Disabled)
- Slack will automatically start/stop based on your configured work schedule

**System Setup:**
- `symlink-setup` - Recreate all configuration file symlinks
- `msys-setup` - Set up MSYS2 development environment

## Repository Structure

- `config/` - All your configuration files:
  - Shell configurations (Bash, Fish, Zsh, PowerShell)
  - Application lists for automatic installation
  - Slack work schedule settings
  - Windows Terminal and Starship prompt settings
  - Symlink mappings for configuration files

- `src/scripts/` - Source code for all the command-line tools

- `src/ps1-windows/` - PowerShell scripts for Windows system configuration (review before running)

- `.build/` - Automatically generated build output (created when you run `__compile.cmd`)

- `__*.cmd` - Setup and installation scripts (run these to get started)

## ⚠️ Important Notes

- **Review Before Running:** The Windows setup scripts (`src/ps1-windows/`) will modify system settings and remove default Windows applications. Please review these scripts before running `__windows-setup.cmd` to ensure they match your preferences.

- **Administrator Rights:** Some scripts require administrator privileges. Windows will prompt you when needed.

- **Backup First:** Consider backing up important data before running system modification scripts.

- **Customization:** All configuration files are in the `config/` directory. Edit these files to customize the setup for your needs, then re-run the appropriate setup scripts.
