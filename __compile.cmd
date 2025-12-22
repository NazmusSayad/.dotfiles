@echo off
setlocal

echo ^> Killing all AHK scripts...
tasklist /NH | findstr /I "^AHK-" >nul && sudo taskkill /F /IM AHK-*

echo ^> Cleaning build directory...
rmdir .\.build\bin /s /q
rmdir .\.build\ahk /s /q
rmdir .\.build\tasks /s /q

echo.
echo ^> Compiling Windows tasks...
call go run ./src/compile-tasks/main.go

echo.
echo ^> Compiling Go scripts...
call go run ./src/compile-go/main.go

echo.
echo ^> Compiling AutoHotkey scripts...
call go run ./src/compile-ahk/main.go
