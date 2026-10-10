---
name: ttyd
description: Opens a terminal in the browser.
---

## Start a server

Run the command below from the directory the terminal should open in, in the background, and note its PID. If port `7681` is taken (check with `lsof -nP -iTCP:7681 -sTCP:LISTEN`), add `-p <free port>`. Do not stop, restart or reuse a `ttyd` server you did not start.

```bash
ttyd -i 127.0.0.1 -W \
  -t fontFamily="FiraCode Nerd Font" -t fontSize=14 -t lineHeight=1.35 -t cursorStyle=bar -t disableLeaveAlert=true \
  -t 'theme={"background":"#1E2227","foreground":"#ABB2BF","cursor":"#FFFFFF","selectionBackground":"#FFFF00","selectionForeground":"#000000","black":"#4A4F5A","red":"#EE5D69","green":"#97C775","yellow":"#B68353","blue":"#61AFEF","magenta":"#C678DD","cyan":"#56B6C2","white":"#CCCCCC","brightBlack":"#616A7B","brightRed":"#FE808A","brightGreen":"#C2DEAD","brightYellow":"#D1A47A","brightBlue":"#93CFFF","brightMagenta":"#E0A8F1","brightCyan":"#82D2DC","brightWhite":"#FFFFFF"}' \
  env -i HOME="$HOME" USER="$USER" LOGNAME="$LOGNAME" SHELL=/opt/homebrew/bin/bash TMPDIR="$TMPDIR" SSH_AUTH_SOCK="$SSH_AUTH_SOCK" \
  LANG=en_US.UTF-8 TERM=xterm-256color COLORTERM=truecolor \
  /opt/homebrew/bin/bash -lc 'exec bash -i' >/tmp/ttyd.log 2>&1 &
```

Open `http://127.0.0.1:<port>`. When done, stop only your server with `kill <PID>`.

## Why the command looks like this

- The theme, font, line height and cursor mirror `~/.dotfiles/config/shell/ghostty.conf`. If that file changes, update the values here to match.
- `env -i` drops the agent's environment so the shell starts like a fresh terminal. It keeps the macOS session variables every native terminal has, because dropping them breaks many tools.
- The login shell loads `profile.sh`, which puts Homebrew and mise on `PATH`. The `exec`'d interactive shell then loads `bashrc.sh`, which sets up starship, zoxide, direnv and mise. Plain `bash` or `bash -i` skips `profile.sh`, so `brew` is missing and the prompt and tools fail to load.
- `/opt/homebrew/bin/bash` is bash 5. `/bin/bash` is macOS's bash 3.2. After `profile.sh`, plain `bash` resolves to bash 5.
- `-i 127.0.0.1` keeps the terminal off the network and `-W` makes it writable. ttyd defaults to all interfaces and read-only.
