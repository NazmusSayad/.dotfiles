@echo off
powershell -Command "if (-not([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] 'Administrator')) { exit 1 }"
if %errorLevel% NEQ 0 (
    echo This script requires administrator privileges.
    echo Press any key to exit...
    pause >nul
    exit /b
)

setlocal
set "SRC=%APPDATA%\Code\User"
set "DST=%APPDATA%\Cursor\User"
set "VS_SRC=%USERPROFILE%\.vscode"
set "CURSOR_DST=%USERPROFILE%\.cursor"

if not exist "%DST%" mkdir "%DST%"

if exist "%DST%\settings.json" del "%DST%\settings.json"
mklink "%DST%\settings.json" "%SRC%\settings.json"

if exist "%DST%\keybindings.json" del "%DST%\keybindings.json"
mklink "%DST%\keybindings.json" "%SRC%\keybindings.json"

if exist "%DST%\snippets" rmdir /S /Q "%DST%\snippets"
mklink /D "%DST%\snippets" "%SRC%\snippets"

if exist "%CURSOR_DST%" rmdir /S /Q "%CURSOR_DST%"
mklink /D "%CURSOR_DST%" "%VS_SRC%"

pause
