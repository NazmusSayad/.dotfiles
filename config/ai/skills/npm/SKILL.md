---
name: npm
description: Run or execute npm scripts and install, update, remove, and manage npm packages. Use this skill when you need to install, update, remove, or manage npm packages. Use this skill whenever you need to run npm scripts (npm run dev, npm run build, etc.) or execute npm packages (npx tsc, npx eslint, etc.).
---

## Commands

- All the commands will be prefixed with `x`.

| Command                     | Description                    | Example                             |
| --------------------------- | ------------------------------ | ----------------------------------- |
| `x ni`                      | Install dependencies           | `x ni`                              |
| `x ni <pkg>`, `ni -D <pkg>` | Add dependency, dev dependency | `x ni react`, `x ni -D types/react` |
| `x nup`, `x nup <pkg>`      | Upgrade dependencies           | `x nup`, `x nup react`              |
| `x nun <pkg>`               | Uninstall dependency           | `x nun <pkg>`                       |
| `x nr <script>`             | Run script                     | `x nr test`                         |
| `x nlx <pkg>`               | Execute package (`npx`)        | `x nlx <package>`                   |

### Query package updates

Run `x taze` to check for available updates for packages.

```sh
x taze
x taze [default|major|minor|patch|latest|newest|next]
```
