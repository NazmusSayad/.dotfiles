---
name: npm
description: Run npm scripts and install, update, remove, and manage npm packages. Use this skill whenever you need to run npm scripts (dev, test, build, etc.). Use this skill when you need to install, update, remove, or manage npm packages.
---

This skill guides npm scripts in real projects and helps run scripts, install packages, update versions, remove packages, and check available updates.

## Run scripts

Use `nr` to run npm scripts (dev, test, build, etc.).

```sh
nr dev --port=3000
nr build
```

## Install packages

Use `ni` to install packages (dependencies, devDependencies, etc.).

```sh
ni <package>
ni -D <package>
```

## Upgrade packages

Use `nup` to upgrade packages (dependencies, devDependencies, etc.).

```sh
nup
nup <package>
```

## Remove packages

Use `nun` to uninstall packages (dependencies, devDependencies, etc.).

```sh
nun <package>
```

## Check available updates

Use `taze` to query package updates (dependencies, devDependencies, etc.).

```sh
taze
taze [default|major|minor|patch|latest|newest|next]
```
