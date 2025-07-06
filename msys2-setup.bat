@echo off
powershell -Command "if (-not([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] 'Administrator')) { exit 1 }"
if %errorLevel% NEQ 0 (
    echo This script requires administrator privileges.
    echo Press any key to exit...
    pause >nul
    exit /b
)

cd /d "%~dp0"

echo Running with administrator privileges.
echo.
echo Running msys-setup...
call npm run msys-setup
echo msys-setup exit code: %errorLevel%

echo.
pause