@echo off
setlocal

call sudo powershell -ExecutionPolicy RemoteSigned -File .\src\etc\install-config.ps1
call proxy-pause echo
