<img width="3840" height="2160" alt="image" src="https://github.com/user-attachments/assets/bac949ae-c913-4727-a77f-ea23a694f02b" />

# .dotfiles

This repository contains reproducible automation and configuration to make it fast and consistent to provision a developer workstation. It includes shell configurations, PowerShell/Batch helpers, and small utilities implemented in Go and TypeScript.

## Structure

- `config/` — Shell configuration and VSCode snippets (e.g. `fish`, VSCode settings).
- `src/` — Source utilities and helper scripts (Go, TypeScript).
- `package.json` — Node tooling and scripts used during development.
- `*.bat`, `*.ps1` — Convenience scripts for installing tools and configuring Windows.

## Common Scripts

- `__install-npm.bat` — Installs Node modules and helper packages referenced in `package.json`.
- `_msys2-setup.bat` — Installs and configures MSYS2 packages.
- `_remove-dev-home.bat` — Cleanup script for removing dev user-specific files.

## Contributing

If you add or modify scripts, follow these guidelines:

- Keep scripts idempotent where possible.
- Document new scripts in this README and add usage examples.
- Test changes locally before opening a pull request.
