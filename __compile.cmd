@echo off

echo.
echo ^> Deleting old compiled files...
del .\bin\*.exe 2>nul

echo.
echo ^> Compiling startup...
call go build -o ".\bin\startup.exe" .\src\startup\main.go

echo.
echo ^> Compiling ahk compiler...
call go build -o ".\bin\ahk-compile.exe" .\src\ahk\main.go

echo.
echo ^> Compiling msys-setup...
call go build -o ".\bin\msys-setup.exe" .\src\msys-setup\main.go

echo.
echo ^> Compiling symlink...
call go build -o ".\bin\symlink-config.exe" .\src\symlink\main.go

echo.
echo ^> Compiling start-menu...
call go build -o ".\bin\start-menu.exe" .\src\start-menu\main.go


echo.
echo ^> Compiling clean-code-snippets...
call go build -o ".\bin\clean-code-snippets.exe" .\src\clean-code-snippets\main.go

echo.
echo ^> Compiling winget-install...
call go build -o ".\bin\winget-install.exe" .\src\winget\install\main.go

echo.
echo ^> Compiling winget-upgrade...
call go build -o ".\bin\winget-upgrade.exe" .\src\winget\upgrade\main.go

echo.
echo ^> Compiling slack-status...
call go build -o ".\bin\slack-status.exe" .\src\slack\status\main.go

echo.
echo ^> Compiling slack-startup...
call go build -o ".\bin\slack-startup.exe" .\src\slack\startup\main.go

echo.
pause
