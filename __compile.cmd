@echo off
setlocal

call sudo powershell -ExecutionPolicy RemoteSigned -File .\src\etc\compile.ps1
call proxy-pause echo
