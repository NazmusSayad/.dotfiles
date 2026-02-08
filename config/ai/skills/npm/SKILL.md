---
name: npm
description: Use @antfu/ni to manage npm packages.
---

This skill guides the use of @antfu/ni to manage npm packages.

The user provides a command to run, install, upgrade, or uninstall a package.

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
