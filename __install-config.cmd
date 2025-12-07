@echo off
setlocal

echo "> Setting up git config..."
git config --global user.name "Nazmus Sayad"
git config --global user.email "87106526+NazmusSayad@users.noreply.github.com"
git config --global init.defaultBranch main
git config --global core.eol lf
git config --global core.autocrlf false
git config --global core.pager cat
git config --global core.ignorecase false
git config --global --add safe.directory "*"
git config --global --add --bool push.autoSetupRemote true

echo "> Installing shell providers..."
pacman -S zsh fish bash

echo "> Installing shell helpers..."
pacman -S mingw-w64-x86_64-starship mingw-w64-clang-x86_64-fastfetch

echo "> Installing helpers..."
pacman -S mingw-w64-x86_64-ffmpeg

echo "> Fish shell config..."
set fish_color_command magenta

echo "> Installing npm packages globally..."
volta install node@latest npm@latest pnpm@latest yarn@latest tsx@latest uni-run@latest code-info@latest netserv@latest

echo "> pnpm config settings..."
pnpm config set ci true
pnpm config set allow-scripts true
pnpm config set shamefully-hoist true
pnpm config set auto-install-peers true

echo.
echo "> Press any key to continue..."
pause >nul
