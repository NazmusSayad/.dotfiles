@echo off
setlocal enabledelayedexpansion
cd /d "%~dp0"

echo ^> Killing all AHK scripts...
set "killed="
for /f "tokens=1" %%p in ('tasklist /fi "IMAGENAME eq AHK-*" /nh 2^>nul ^| findstr /i "AHK-"') do (
    set "name=%%p"
    set "name=!name:.exe=!"
    set "killed=!killed! !name!"
)
if defined killed sudo taskkill /F /IM AHK-* >nul 2>&1

echo ^> Cleaning AHK build directory...
if exist ".build\ahk" rmdir /s /q ".build\ahk"

echo.
echo ^> Compiling AutoHotkey scripts...
call go run ./src/compile-ahk/main.go

if defined killed (
    echo.
    echo ^> Restarting AHK scripts...
    for %%n in (%killed%) do (
        set "exe=%CD%\.build\ahk\%%n.exe"
        if exist "!exe!" (
            echo ^> Restarting: !exe!
            powershell -NoProfile -Command "Start-Process '!exe!' -Verb RunAs"
        )
    )
)

echo.
echo Done!
