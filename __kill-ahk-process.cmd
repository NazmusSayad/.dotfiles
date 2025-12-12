@echo off
setlocal

echo ^> Killing all AHK scripts...
sudo taskkill /F /IM AHK-*

echo.
echo "> Press any key to continue..."
pause >nul