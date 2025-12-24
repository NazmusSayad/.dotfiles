@echo off
setlocal

echo ^> Setting up Go environment variables...
powershell.exe -NoProfile -Command "[Environment]::SetEnvironmentVariable('GOPATH', (go env GOPATH), 'User')"
powershell.exe -NoProfile -Command "[Environment]::SetEnvironmentVariable('GOROOT', (go env GOROOT), 'User')"

echo ^> Installing Go tools...
call go install -v golang.org/x/tools/gopls@latest
call go install -v mvdan.cc/sh/v3/cmd/shfmt@latest
call go install -v honnef.co/go/tools/cmd/staticcheck@latest

echo.
pause