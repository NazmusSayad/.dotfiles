@echo off

call go build -o ./build/clean-vscode-snippets.exe ./src/clean-code-snippets.go

echo.
pause