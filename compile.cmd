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
echo ^> Compiling clean-code-snippets...
call go build -o .\___clean-vscode-snippets.exe .\src\clean-code-snippets\main.go

echo.
echo ^> Compiling fish-greetings...
call go build -o .\___fish-greetings.exe .\src\fish-greetings\main.go

echo.
echo ^> Compiling winget-install...
call go build -o .\___winget-install.exe .\src\winget\install\main.go

echo.
echo ^> Compiling winget-upgrade...
call go build -o .\___winget-upgrade.exe .\src\winget\upgrade\main.go

echo.
echo ^> Compiling winget-upgrade-auto...
call go build -o .\___winget-upgrade-auto.exe .\src\winget\upgrade-auto\main.go

echo.
pause
