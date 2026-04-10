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

## Rules

- For dependency and script operations, use the commands mentioned above.
- Never use `npm`, `pnpm`, `yarn`, or `bun` unless explicitly instructed.
- NEVER UPGRADE **major** VERSIONS of dependencies without explicit instructions. Always ask for confirmation before doing so.

## Guidelines

- Use `taze` before upgrading dependencies to check for available updates.
- `taze` is not an npm package, it's a standalone tool; do not run it via `npx` or `nlx`.
