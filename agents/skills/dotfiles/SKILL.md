---
name: dotfiles
description: The user's rules for tools and environment variables. Use when you install a tool, change a tool version, or set environment variables.
---

## Package/Tool Managers

**macOS**: Use mise. If mise does not have the tool, use `brew`. Ask the user before you use `brew`.

**Linux**: Use mise. If mise does not have the tool, use `snap`. If `snap` does not have the tool, use `apt`. Ask the user before you use `snap` or `apt`.

## Local first

If the tool is not in the global mise config, install it in the project with mise. If the project needs a different version than the global one, set that version in the project with mise. Do not change the global mise config unless the user tells you to.

## Environments

Use the environment files that the framework supports, for example `.env` or `.env.local`. If the variables must be in the shell, use `.envrc` with direnv.
