---
name: npm
description: Install, update, remove and manage npm packages; run scripts and execute packages via npx. Must use this skill instead of npm, pnpm, yarn, or bun for dependency and workflow management.
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

## Guidelines

- Never use `npm`, `pnpm`, `yarn`, or `bun` unless explicitly instructed.
- For dependency and script operations, always use the command set in this file.
- `taze` is not an npm package; do not run it via `npx` or `nlx`.
- NEVER UPGRADE **major** VERSIONS of dependencies without explicit instructions. Always ask for confirmation before doing so.
