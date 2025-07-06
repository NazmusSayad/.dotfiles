@echo off

echo Removing Windows Dev Home...
powershell -Command "Get-AppxPackage -AllUsers -PackageTypeFilter Bundle -Name '*Windows.DevHome*' | Remove-AppxPackage -AllUsers"

echo.
echo Dev Home removal completed.
echo Press any key to continue...
pause >nul
