@echo off
powershell -Command "if (-not([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] 'Administrator')) { exit 1 }"
if %errorLevel% NEQ 0 (
    echo This script requires administrator privileges.
    echo Press any key to exit...
    pause >nul
    exit /b
)

cd /d "%~dp0"
echo Upgrading packages from winget-apps.txt...

for /f "usebackq delims=" %%i in ("winget-apps.txt") do (
    echo Upgrading: %%i
    winget upgrade "%%i" --interactive --no-upgrade --accept-package-agreements --accept-source-agreements
)

echo Done!
pause
