@echo off

echo.
echo ^> Compiling ahk compiler...
call go build -o .\___ahk-compile.exe .\src\ahk\main.go

echo.
echo ^> Compiling ahk scripts...
call .\___ahk-compile.exe

echo.
echo ^> Compiling msys-setup...
call go build -o .\___msys-setup.exe .\src\msys-setup\main.go

echo.
echo ^> Compiling symlink...
call go build -o .\___symlink-config.exe .\src\symlink\main.go

echo.
echo ^> Compiling clean-code-snippets...
call go build -o .\___clean-vscode-snippets.exe .\src\clean-code-snippets\main.go

echo.
echo ^> Compiling winget-install...
call go build -o .\___winget-install.exe .\src\winget\install\main.go

echo.
echo ^> Compiling winget-upgrade...
call go build -o .\___winget-upgrade.exe .\src\winget\upgrade\main.go

echo.
echo ^> Compiling slack-enable...
call go build -o .\___slack-enable.exe .\src\slack-enable\main.go

echo.
echo ^> Compiling slack-disable...
call go build -o .\___slack-disable.exe .\src\slack-disable\main.go

echo.
pause
