@echo off
setlocal

echo ^> Setting up Go environment variables...
for /f "usebackq delims=" %%G in (`go env GOPATH`) do setx GOPATH "%%G"
for /f "usebackq delims=" %%G in (`go env GOROOT`) do setx GOROOT "%%G"

echo ^> Installing Go tools...
call go install -v golang.org/x/tools/gopls@latest
call go install -v mvdan.cc/sh/v3/cmd/shfmt@latest
call go install -v honnef.co/go/tools/cmd/staticcheck@latest

echo.
pause