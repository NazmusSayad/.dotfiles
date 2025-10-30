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
set "DST_1=%APPDATA%\Cursor\User"
set "DST_2=%APPDATA%\Windsurf\User"


if not exist "%DST_1%" mkdir "%DST_1%"
if not exist "%DST_2%" mkdir "%DST_2%"

if exist "%DST_1%\settings.json" del "%DST_1%\settings.json"
mklink "%DST_1%\settings.json" "%SRC%\settings.json"

if exist "%DST_2%\settings.json" del "%DST_2%\settings.json"
mklink "%DST_2%\settings.json" "%SRC%\settings.json"

if exist "%DST_1%\keybindings.json" del "%DST_1%\keybindings.json"
mklink "%DST_1%\keybindings.json" "%SRC%\keybindings.json"

if exist "%DST_2%\keybindings.json" del "%DST_2%\keybindings.json"
mklink "%DST_2%\keybindings.json" "%SRC%\keybindings.json"

if exist "%DST_1%\snippets" rmdir /S /Q "%DST_1%\snippets"
mklink /D "%DST_1%\snippets" "%SRC%\snippets"

if exist "%DST_2%\snippets" rmdir /S /Q "%DST_2%\snippets"
mklink /D "%DST_2%\snippets" "%SRC%\snippets"

pause
