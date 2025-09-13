@echo off

echo.
echo Compiling ahk scripts
call npm run ahk-compile

echo.
echo Compiling clean-code-snippets
call go build -o ./___winget-install.exe ./src/winget/install/main.go
call go build -o ./___winget-upgrade.exe ./src/winget/upgrade/main.go
call go build -o ./___clean-vscode-snippets.exe ./src/clean-code-snippets/main.go

echo.
echo.
pause