@echo off
setlocal enabledelayedexpansion

echo ^> Killing all AHK scripts...
sudo taskkill /F /IM AHK-*

echo ^> Cleaning build directory...
rmdir .\build\bin /s /q
rmdir .\build\ahk /s /q

echo ^> Compiling Go scripts...
go run ./src/compile-scripts/main.go

echo.
echo ^> Compiling Rust functions...
go run ./src/compile-functions/main.go

echo.
echo ^> Compiling AutoHotkey scripts...
go run ./src/compile-ahk/main.go

echo.
pause