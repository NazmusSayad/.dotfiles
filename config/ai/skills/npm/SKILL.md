---
name: npm
description: Install, update, remove and manage npm packages; run scripts and execute packages via npx. Use this skill instead of npm, pnpm, yarn, or bun for dependency and workflow management.
---

## Commands

| Command                       | Description                   | Example                             |
| ----------------------------- | ----------------------------- | ----------------------------------- |
| `x ni`                        | Install dependencies          | `x ni`                              |
| `x ni <pkg>`, `x ni -D <pkg>` | Add dependency, devDependency | `x ni react`, `x ni -D types/react` |
| `x nup`, `x nup <pkg>`        | Upgrade dependencies          | `x nup`, `x nup react`              |
| `x nun <pkg>`                 | Uninstall dependency          | `x nun <pkg>`                       |
| `x nr <script>`               | Run script                    | `x nr test`                         |
| `x nlx <pkg>`                 | Execute package (`npx`)       | `x nlx <package>`                   |
| `x taze`                      | Check package updates         | `x taze major`                      |

**NOTE:** All the commands will be prefixed with `x`.
