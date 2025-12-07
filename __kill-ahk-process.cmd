@echo off
setlocal

net session >nul 2>&1
if %errorLevel% NEQ 0 (
    echo FAIL: Administrator privileges required.
    echo Relaunching with elevated privileges...
    powershell -NoProfile -Command "Start-Process -FilePath '%~f0' -Verb RunAs"
    exit /b
)

echo ^> Killing all AHK scripts...
taskkill /F /IM AHK-*