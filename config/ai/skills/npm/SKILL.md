---
name: npm
description: Installs, upgrades, and removes Node.js/npm dependencies, runs scripts, executes packages, and manages package workflows. MUST USE this skill for any task involving npm, npx, pnpm, yarn, or bun, including package.json dependency changes, lockfile updates, install/update/uninstall actions, and script execution.
---

## Commands

| Command                   | Description                   | Example                          |
| ------------------------- | ----------------------------- | -------------------------------- |
| `ni`                      | Install dependencies          | `ni`                             |
| `ni <pkg>`, `ni -D <pkg>` | Add dependency, devDependency | `ni react`, `ni -D @types/react` |
| `nup`, `nup <pkg>`        | Upgrade dependencies          | `nup`, `nup react`               |
| `nun <pkg>`               | Uninstall dependency          | `nun react`                      |
| `nr <script>`             | Run npm script                | `nr test`                        |
| `nlx <pkg>`               | Execute npm package (`npx`)   | `nlx eslint`                     |
| `taze`                    | Check package updates         | `taze major`                     |

## Rules

- For dependency and script operations, use the commands mentioned above.
- NEVER use `npm`, `pnpm`, `yarn`, or `bun` directly unless explicitly instructed.
- NEVER upgrade **major** versions of dependencies without explicit instructions. Always ask for confirmation before doing so.

## Guidelines

- Use `taze` before upgrading dependencies to check for available updates.
- `taze` is not an npm package, it's a standalone tool; do not run it via `npx` or `nlx`.
