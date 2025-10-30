@echo off
powershell -Command "if (-not([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] 'Administrator')) { exit 1 }"
if %errorLevel% NEQ 0 (
    echo This script requires administrator privileges.
    echo Press any key to exit...
    pause >nul
    exit /b
)

setlocal
set SRC=%~dp0\config\zed
set "DEST=%APPDATA%\Zed"

if exist "%DEST%" rmdir /S /Q "%DEST%"
mklink /D "%DEST%" "%SRC%"

pause
