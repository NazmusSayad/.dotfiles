---
name: npm
description: Run npm scripts and install, update, remove, manage npm packages. Use this skill whenever you need to run npm scripts (dev, test, build, etc.) or manage packages.
---

This skill guides npm scripts in real projects and helps run scripts, add packages, update versions, remove packages, and check available updates.

## Run scripts

Use `nr` to run npm scripts.

```bash
nr dev --port=3000
nr build
```

## Install packages

Use `ni` to install packages.

```bash
ni <package>
ni -D <package>
```

## Upgrade packages

Use `nup` to upgrade dependencies.

```bash
nup
nup <package>
```

## Remove packages

Use `nun` to uninstall packages.

```bash
nun <package>
```

## Check available updates

Use `taze` to query dependency updates.

```bash
taze
taze [default|major|minor|patch|latest|newest|next]
```
