---
name: npm
description: Use @antfu/ni to run scripts, install/upgrade/uninstall packages across npm, yarn, and pnpm. Use when running npm/yarn/pnpm commands, installing dependencies, or managing packages.
---

Use **@antfu/ni** (ni/nr/nup/nun). It picks the project's package manager automatically.

## `nr` - run

```bash
nr dev --port=3000
nr build
```

## `ni` - install

```bash
ni <package>
ni -D <package>
```

## `nup` - upgrade

```bash
nup
nup <package>
```

## `nun` - uninstall

```bash
nun <package>
```
