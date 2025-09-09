@echo off
setlocal enabledelayedexpansion

powershell -Command "if (-not([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] 'Administrator')) { exit 1 }"
if %errorLevel% NEQ 0 (
    echo This script requires administrator privileges.
    echo Press any key to exit...
    pause >nul
    exit /b
)

cd /d "%~dp0"
@REM START OF SCRIPT

set PACKAGES=

for /f "usebackq delims=" %%a in ("winget-apps.ini") do (
  set line=%%a
  if not "!line!"=="" if not "!line:~0,1!"=="#" (
    for /f "tokens=1" %%b in ("!line!") do (
      set package=%%b
      if not "!package!"=="" set PACKAGES=!PACKAGES! !package!
    )
  )
)

for /f "tokens=* delims= " %%a in ("!PACKAGES!") do set PACKAGES=%%a

if not "!PACKAGES!"=="" (
  winget upgrade !PACKAGES! --no-upgrade --interactive --accept-package-agreements --accept-source-agreements
) else (
  echo No valid packages found
)

pause
