@echo off
setlocal

echo ^> Killing all AHK scripts...
tasklist /NH | findstr /I "^AHK-" >nul && sudo taskkill /F /IM AHK-*

echo ^> Cleaning build directory...
rmdir .\.build\bin /s /q >nul
rmdir .\.build\ahk /s /q >nul

echo.
echo ^> Compiling AutoHotkey scripts...
call go run ./src/compile-ahk/main.go

echo.
echo ^> Compiling aliases...
call go run ./src/compile-alias/main.go

echo.
echo ^> Compiling Go scripts...
call go run ./src/compile-scripts/main.go
