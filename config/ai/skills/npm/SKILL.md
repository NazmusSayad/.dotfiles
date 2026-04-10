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

- Never use `npm`, `pnpm`, `yarn`, or `bun` directly for dependency or workflow management. Always use the commands listed above.
- NEVER upgrade **major** versions of dependencies without explicit instructions. Always ask for confirmation before doing so.
- `taze` is not an npm package. Never run it with `nlx` or `npx`.
