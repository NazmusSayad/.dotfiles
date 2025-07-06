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

setlocal enabledelayedexpansion

:: Set execution policy equivalent (batch files run by default)
:: Execute all PowerShell scripts in the scripts directory
for %%f in (.\src\ps1\*.ps1) do (
    echo Executing: %%f
    powershell.exe -ExecutionPolicy RemoteSigned -File "%%f"
    if !errorlevel! neq 0 (
        echo Error executing %%f
        pause
        exit /b !errorlevel!
    )
)

:: Restart computer
echo.
echo All scripts executed successfully.
echo.
echo Press any key to restart the computer...
pause >nul
echo Restarting computer...
shutdown /r /f /t 0
