---
name: npm
description: Installs, upgrades, and removes Node.js/npm dependencies, runs scripts, executes packages, and manages package workflows. MUST USE this skill for any task involving npm, npx, pnpm, yarn, or bun, including package.json dependency changes, lockfile updates, install/update/uninstall actions, and script execution.
---

## Commands

| Command                         | Description                       | Example                                |
| ------------------------------- | --------------------------------- | -------------------------------------- |
| `ni`, `ni <pkg>`, `ni -D <pkg>` | Install dependency, devDependency | `ni`, `ni react`, `ni -D @types/react` |
| `nup`, `nup <pkg>`              | Upgrade dependencies              | `nup`, `nup react`                     |
| `nun <pkg>`                     | Uninstall dependency              | `nun react`                            |
| `nr <script>`                   | Run npm script                    | `nr lint`, `nr test`                   |
| `nlx <pkg>`                     | Execute npm package (`npx`)       | `nlx tsc`                              |

## Rules

- NEVER use `npm`, `pnpm`, `yarn`, or `bun` directly unless explicitly instructed. USE the commands mentioned above.
- Prefer existing `package.json` scripts via `nr <script>` (e.g. `lint`, `typecheck`, `format`) over running tools directly with `nlx`/`npx` when available. For example, use `nr lint` instead of `nlx eslint` or `nr typecheck` instead of `nlx tsc`. Use `nlx <tool>` when `package.json` has no matching script.
