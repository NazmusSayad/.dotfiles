@echo off
setlocal

sudo powershell -ExecutionPolicy RemoteSigned -File .\src\others\compile.ps1

echo.
pause