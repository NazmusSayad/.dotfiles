@echo off

echo.
echo ^> Deleting old compiled files...
del /Q .\___*.exe 2>nul

echo.
echo ^> Compiling launch...
call go build -o ".\___Launch.exe" .\src\launch\main.go

echo.
echo ^> Compiling ahk compiler...
call go build -o ".\___AHK Compile.exe" .\src\ahk\main.go

echo.
echo ^> Compiling ahk scripts...
call ".\___AHK Compile.exe"

echo.
echo ^> Compiling msys-setup...
call go build -o ".\___MSYS Setup.exe" .\src\msys-setup\main.go

echo.
echo ^> Compiling symlink...
call go build -o ".\___Symlink Config.exe" .\src\symlink\main.go

echo.
echo ^> Compiling clean-code-snippets...
call go build -o ".\___Clean VSCode Snippets.exe" .\src\clean-code-snippets\main.go

echo.
echo ^> Compiling winget-install...
call go build -o ".\___Winget Install.exe" .\src\winget\install\main.go

echo.
echo ^> Compiling winget-upgrade...
call go build -o ".\___Winget Upgrade.exe" .\src\winget\upgrade\main.go

echo.
echo ^> Compiling slack-status...
call go build -o ".\___Slack Status.exe" .\src\slack-status\main.go

echo.
echo ^> Compiling slack-launch...
call go build -o ".\___Slack Launch.exe" .\src\slack-launch\main.go

echo.
pause
