@echo off
setlocal

echo Setting up GOPATH...
call powershell.exe -NoProfile -Command "[Environment]::SetEnvironmentVariable('GOPATH', (go env GOPATH), 'User')"

echo Setting up GOROOT...
call powershell.exe -NoProfile -Command "[Environment]::SetEnvironmentVariable('GOROOT', (go env GOROOT), 'User')"

echo.
echo Done! Press any key to exit...
pause>nul
