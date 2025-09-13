@echo off

echo.
echo Compiling ahk scripts
call npm run ahk-compile

echo.
echo Compiling clean-code-snippets
call go build -o ./___clean-vscode-snippets.exe ./src/clean-code-snippets/main.go

echo.
echo.
pause