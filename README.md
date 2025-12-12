<img width="2855" alt="image" src="https://github.com/user-attachments/assets/f217bd38-32a5-4bee-a259-96d3f5fb6837" />

# Development Setup for Windows

This repository contains reproducible automation and configuration to make it fast and consistent to provision a developer workstation. It includes shell configurations, PowerShell/Batch helpers, and small utilities implemented in Go.

## Features & Capabilities

- **Windows Configuration**
  Enhances productivity with custom AutoHotkey scripts for keyboard-centric window management, rapid virtual desktop switching, and tailored input configurations.

- **Apps, Packages, and Runtimes Management**
  Automates the lifecycle of software dependencies. Orchestrates installation and updates for system applications and global runtimes, ensuring a reproducible development stack.

- **Shell Experience**
  Delivers a consistent, Unix-like environment on Windows. Fully configured setups for Bash and Fish include cross-shell prompts, aliases, and modern CLI tools.

- **Code Editor Configuration**
  Synchronizes Visual Studio Code preferences, keybindings, and extensions. Includes dedicated helpers to manage code snippets and maintain a unified editing environment.

- **Communication Optimization**
  Streamlines workspace connectivity with automated tools for managing application states. Includes intelligent Slack launch management with work-hour scheduling (automatically starts/stops Slack based on Bangladesh time zone), and AHK process management utilities.

- **System Performance**
  Maximizes hardware potential by debloating Windows. Scripts aggressively remove unused pre-installed apps, disable telemetry services, and tune system policies for development workloads.

## Getting Started

### Prerequisites

- Windows 10 or Windows 11
- Git (installed to clone the repository)
- Go (for compiling the utilities)

### Installation Guide

1.  **Clone the Repository:**

    ```shell
    git clone https://github.com/NazmusSayad/.dotfiles.git
    ```

2.  **Initial Setup:**
    Run the `__install.cmd` script as Administrator. This script will:

    - Configure global Git settings.
    - Install the Volta package manager.
    - Install essential global npm packages (Node, pnpm, yarn, etc.).

3.  **Build Utilities:**
    Run the `__compile.cmd` script. This uses Go to compile the helper utilities located in `src/` into the `build/` directory.

4.  **Install Dotfiles:**
    Run the `__install-dotfiles.cmd` script as Administrator. This creates a symbolic link from `%USERPROFILE%\.dotfiles` to your repository and adds the `build/` directory to your system PATH.

5.  **Apply Configurations:**
    - **Symlinks:** Run `symlink-config.exe` to link your config files.
    - **Software:** Run `winget-install.exe` to install applications defined in `config/winget-apps.jsonc`.
    - **System Settings:** Review and run the PowerShell scripts in `src/ps1/` as needed (e.g., `settings.ps1`) to apply system optimizations.
    - **Slack Management:** Use `slack-status.exe` to configure intelligent Slack launch behavior (Always, Work Hours, or Disabled).

## Repository Structure

- `build/`: Compiled Go executables (gitignored).
- `config/`: Configuration files for shells, standard apps, and `winget` packages.
- `src/`: Source code for Go utilities and PowerShell scripts.
  - `src/ps1/`: Windows debloating and configuration scripts.
  - `src/ahk/`: AutoHotkey scripts for window management.
  - `src/winget/`: Tools for parsing and installing Winget packages.
  - `src/slack/`: Slack automation utilities with intelligent launch scheduling.
- `bin/`: Shared binaries.
- `__*`: Installation and utility scripts.

## ⚠️ Disclaimer

This repository contains scripts that modify system settings and remove default applications. Review all scripts (especially those in `src/ps1/`) before running them to ensure they align with your requirements.
