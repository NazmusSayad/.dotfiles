@echo off
setlocal EnableDelayedExpansion

set "CURRENT_DIR=%~dp0"
set "CURRENT_DIR=%CURRENT_DIR:~0,-1%"

set "DOTFILES_DIR=%CURRENT_DIR%"
set "DOTFILES_DIR_BIN=%DOTFILES_DIR%\.build\bin"

echo Setting DOTFILES_DIR...
powershell -NoProfile -Command "[Environment]::SetEnvironmentVariable('DOTFILES_DIR','%DOTFILES_DIR%','User')"

set "USER_PATH="
for /f "usebackq delims=" %%A in (`powershell -NoProfile -Command "[Environment]::GetEnvironmentVariable('PATH','User')"`) do set "USER_PATH=%%A"

echo %USER_PATH% | find /I "%DOTFILES_DIR_BIN%" >nul 2>&1
if errorlevel 1 (
  echo DOTFILES_DIR_BIN not in PATH
  echo Adding DOTFILES_DIR_BIN to PATH...
  powershell -NoProfile -Command "[Environment]::SetEnvironmentVariable('PATH','%DOTFILES_DIR_BIN%;%USER_PATH%','User')"
) else (
  echo DOTFILES_DIR_BIN already in PATH
)

echo.
pause