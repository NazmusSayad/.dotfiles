@echo off
powershell -Command "if (-not([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] 'Administrator')) { exit 1 }"
if %errorLevel% NEQ 0 (
    echo This script requires administrator privileges.
    echo Press any key to exit...
    pause >nul
    exit /b
)

setlocal EnableDelayedExpansion

set SRC=%~dp0
set "DEST=%APPDATA%\Microsoft\Windows\Start Menu\Programs\dotfiles"

if exist "%DEST%" rmdir /S /Q "%DEST%"
mkdir "%DEST%"

for %%F in ("%SRC%*") do (
    if not exist "%%~fF\" (
        set "n=%%~nxF"
        set "name="
        echo Considering file: %%~nxF
        if "!n:~0,3!"=="___" set "name=!n:~3!"
        if not defined name if "!n:~0,2!"=="__" set "name=!n:~2!"
        if defined name (
            echo Creating shortcut: !name!.lnk -> %%~nxF
            powershell -NoProfile -ExecutionPolicy Bypass -Command "$w=New-Object -ComObject WScript.Shell; $s=$w.CreateShortcut('%DEST%\!name!.lnk'); $s.TargetPath='%%~fF'; $s.WorkingDirectory='%%~dpF'; $s.IconLocation='%%~fF,0'; $s.Save()" >nul
            if !errorlevel! EQU 0 (
                echo OK: !name!.lnk
            ) else (
                echo FAIL: !name!.lnk
            )
        ) else (
            echo Skipped: %%~nxF
        )
    ) else (
        echo Skipped directory: %%~nxF
    )
)

pause
