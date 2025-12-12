@echo off
setlocal

echo.
echo ^> Installing start menu entries...
call go run ./src/install-start-menu/main.go

echo.
pause