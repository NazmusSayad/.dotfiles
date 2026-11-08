---
name: fleet
description: Use when a task needs to reach or run something on another of the user's machines.
---

## Omnitrix

- OS: macOS
- SSH: `ssh <redacted> 'bash -lc "<command>"'`

## Ultimatrix

- OS: Windows
- SSH: `ssh <redacted>` (runs PowerShell 7)
- WSL: `ssh <redacted> 'bash -lc "<command>"'`

> Note: Development happens primarily in WSL (Ubuntu 26.04).

> Note: SSH goes through Cloudflare Access. If `ssh` prints a Cloudflare login URL and waits, ask the user to log in.
