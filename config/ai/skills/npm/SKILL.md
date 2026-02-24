---
name: npm
description: Run or execute npm scripts and install, update, remove, and manage npm packages. Use this skill whenever you need to run npm scripts (npm run dev, npm run build, etc.) or execute npm packages (npx tsc, npx eslint, etc.). Use this skill when you need to install, update, remove, or manage npm packages.
---

This skill guides npm scripts in real projects and helps run scripts, install packages, update versions, remove packages, and check available updates.

## Run scripts

Run `nr` to run npm scripts (dev, test, start, build, etc.).

```sh
nr <script>
nr dev --port=3000
nr build
```

## Execute packages

Run `nlx` to execute npm packages (tsc, eslint, etc.).

```sh
nlx <package>
nlx next dev
```

## Install packages

Run `ni` to install packages (dependencies, devDependencies, etc.).

```sh
ni <package>
ni -D <package>
```

## Upgrade packages

Run `nup` to upgrade packages (dependencies, devDependencies, etc.).

```sh
nup
nup <package>
```

### Query package updates

Run `taze` to check for available updates for packages.

```sh
taze
taze [default|major|minor|patch|latest|newest|next]
```

## Remove packages

Run `nun` to uninstall packages (dependencies, devDependencies, etc.).

```sh
nun <package>
```
