@echo off
setlocal

set "CURRENT_DIR=%~dp0"
set "CURRENT_DIR=%CURRENT_DIR:~0,-1%"

set "DOTFILES_DIR=%CURRENT_DIR%"
set "DOTFILES_DIR_BIN=%DOTFILES_DIR%\.build\bin"

echo Setting DOTFILES_DIR...
powershell -NoProfile -Command "[Environment]::SetEnvironmentVariable('DOTFILES_DIR','%DOTFILES_DIR%','User')"

