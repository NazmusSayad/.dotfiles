---
name: fleet
description: Use when a task needs to reach or run something on another of the user's machines.
---

## Omnitrix

- OS: macOS
- SSH: `ssh 192.168.0.100 'bash -lc "<command>"'`

## Ultimatrix

- OS: Windows
- SSH: `ssh 192.168.0.102` (runs PowerShell 7)
- WSL: `ssh 192.168.0.102 'wsl.exe -d Ubuntu-26.04 -u sayad --cd "~" -- bash -lc "<command>"'`

> Note: Development happens primarily in WSL (Ubuntu 26.04).
