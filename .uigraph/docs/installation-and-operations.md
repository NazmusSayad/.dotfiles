# Installation and operations

## Supported environments

The repository provisions Windows 10/11 and macOS workstations. Git is required. Go 1.24.5 is declared in `go.mod`; MSYS2 and Scoop are used by the Windows setup, and Homebrew is used by the macOS setup.

## Setup commands

- Windows: run `__install-dotfiles.cmd` to set the user `DOTFILES_DIR` variable and prepend `.build/bin` to the user PATH. Run `__install-config.cmd` to apply configuration links and platform tasks.
- Windows system changes: run `__setup-windows.cmd` as Administrator. It executes every `src/ps1-windows/*.ps1` script and restarts the machine.
- macOS: review and run `__setup-macos.sh`. It sources `etc/git-config.sh` and `etc/shell-config.sh`, changes shell and macOS defaults, opens the Chrome profile, and restarts Dock, Finder, and SystemUIServer.
- Build utilities: `__compile.sh` removes `.build/bin` and invokes `go run ./src/compile-scripts/main.go`. Repository guidance says not to run build or setup commands automatically.

## Privileges and side effects

Windows system setup requires Administrator privileges. The macOS setup uses `sudo` for system settings, shell registration, battery controls, and service settings. Setup writes global Git configuration, modifies user environment variables, creates configuration links, configures startup tasks or launch agents, and may disable or alter operating-system services. Review scripts and back up the workstation before use.

## Operational boundaries

The toolkit invokes OS commands including `git`, `bash`, `fish`, `defaults`, `launchctl`, `mdutil`, `chflags`, `chsh`, `open`, `killall`, PowerShell, `shutdown`, and package managers. The GitHub release helper calls `https://api.github.com/repos/{owner}/{repo}/releases/latest`, then downloads the selected asset URL; HTTP failures are returned and callers may terminate the utility.
