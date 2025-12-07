@echo off

echo ^> Deleting old compiled files...
del .\build /s /q

echo.
echo ^> Compiling windows startup...
call go build -o ".\build\windows-startup.exe" .\src\startup\main.go

echo.
echo ^> Compiling ahk compiler...
call go build -o ".\build\ahk-compile.exe" .\src\ahk\main.go

echo.
echo ^> Compiling msys-setup...
call go build -o ".\build\msys-setup.exe" .\src\msys-setup\main.go

echo.
echo ^> Compiling symlink...
call go build -o ".\build\symlink-config.exe" .\src\symlink\main.go

echo.
echo ^> Compiling clean-code-snippets...
call go build -o ".\build\clean-code-snippets.exe" .\src\clean-code-snippets\main.go

echo.
echo ^> Compiling winget-install...
call go build -o ".\build\winget-install.exe" .\src\winget\install\main.go

echo.
echo ^> Compiling winget-upgrade...
call go build -o ".\build\winget-upgrade.exe" .\src\winget\upgrade\main.go

echo.
echo ^> Compiling slack-status...
call go build -o ".\build\slack-status.exe" .\src\slack\status\main.go

echo.
echo ^> Compiling slack-startup...
call go build -o ".\build\slack-startup.exe" .\src\slack\startup\main.go

echo.
echo ^> Compiling ahk scripts...
call .\build\ahk-compile.exe

echo.
pause
