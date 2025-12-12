@echo off
setlocal

echo ^> Installing go...
call go install -v golang.org/x/tools/gopls@latest
call go install -v mvdan.cc/sh/v3/cmd/shfmt@latest
call go install -v honnef.co/go/tools/cmd/staticcheck@latest
