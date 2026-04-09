---
name: npm
description: Install, update, remove and manage npm packages; run scripts and execute packages via npx. Must use this skill instead of npm, pnpm, yarn, or bun for dependency and workflow management.
---

## Commands

| Command                   | Description                   | Example                         |
| ------------------------- | ----------------------------- | ------------------------------- |
| `ni`                      | Install dependencies          | `ni`                            |
| `ni <pkg>`, `ni -D <pkg>` | Add dependency, devDependency | `ni react`, `ni -D types/react` |
| `nup`, `nup <pkg>`        | Upgrade dependencies          | `nup`, `nup react`              |
| `nun <pkg>`               | Uninstall dependency          | `nun <pkg>`                     |
| `nr <script>`             | Run script                    | `nr test`                       |
| `nlx <pkg>`               | Execute package (`npx`)       | `nlx <package>`                 |
| `taze`                    | Check package updates         | `taze major`                    |
