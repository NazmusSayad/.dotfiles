---
name: npm
description: Installs, upgrades, and removes Node.js/npm dependencies, runs scripts, executes packages, and manages package workflows. MUST USE for any task involving npm, npx, pnpm, yarn, or bun, including package.json dependency changes, lockfile updates, install/update/uninstall actions, and script execution.
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

- NEVER use `npm`, `pnpm`, `yarn`, or `bun` directly unless explicitly instructed by the user. USE the commands mentioned above.
- NEVER run a tool directly with `nlx`/`npx` when `package.json` has a matching script (e.g. `lint`, `typecheck`, `format`). USE `nr <script>` instead. For example, `nr lint` instead of `nlx eslint`, or `nr typecheck` instead of `nlx tsc`. Use `nlx <tool>` only when no matching script exists.
- When the project's own instructions explicitly give a command to run, run that command instead of this skill's equivalent.
