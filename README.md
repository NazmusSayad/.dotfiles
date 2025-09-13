# .dotfiles

This repository contains reproducible automation and configuration to make it fast and consistent to provision a developer workstation. It includes shell configurations, PowerShell/Batch helpers, and small utilities implemented in Go and TypeScript.

## Structure

- `config/` — Shell configuration and VSCode snippets (e.g. `fish`, VSCode settings).
- `src/` — Source utilities and helper scripts (Go, TypeScript).
- `package.json` — Node tooling and scripts used during development.
- `*.bat`, `*.ps1` — Convenience scripts for installing tools and configuring Windows.

## Quick Start

1. Inspect the repository and decide which scripts you want to run.
2. Run interactive setup scripts from an elevated PowerShell when required, for example:

   - `__install-npm.bat` — installs Node and related packages used by the repo.
   - `_msys2-setup.bat` — bootstraps MSYS2 when you need a POSIX-like toolchain.

3. Apply configuration files from `config/` to your user profile or review them before running.

## Common Scripts

- `__install-npm.bat` — Installs Node modules and helper packages referenced in `package.json`.
- `_msys2-setup.bat` — Installs and configures MSYS2 packages.
- `_remove-dev-home.bat` — Cleanup script for removing dev user-specific files.

## Contributing

If you add or modify scripts, follow these guidelines:

- Keep scripts idempotent where possible.
- Document new scripts in this README and add usage examples.
- Test changes locally before opening a pull request.
