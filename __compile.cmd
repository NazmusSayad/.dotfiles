@echo off

echo ^> Deleting old compiled files...
del .\build\bin /s /q

echo.
echo ^> Compiling windows startup (Go)...
go build -o ".\build\bin\windows-startup.exe" .\src\startup\main.go

echo.
echo ^> Compiling vbproxy (Go)...
go build -o ".\build\bin\vbproxy.exe" .\src\vbproxy\main.go

echo.
echo ^> Compiling ahk compiler (Go)...
go build -o ".\build\bin\ahk-compile.exe" .\src\ahk\main.go

echo.
echo ^> Compiling msys-setup (Go)...
go build -o ".\build\bin\msys-setup.exe" .\src\msys-setup\main.go

echo.
echo ^> Compiling symlink (Go)...
go build -o ".\build\bin\symlink-config.exe" .\src\symlink\main.go

echo.
echo ^> Compiling clean-code-snippets (Go)...
go build -o ".\build\bin\clean-code-snippets.exe" .\src\clean-code-snippets\main.go

echo.
echo ^> Compiling winget-install (Go)...
go build -o ".\build\bin\winget-install.exe" .\src\winget\install\main.go

echo.
echo ^> Compiling winget-upgrade (Go)...
go build -o ".\build\bin\winget-upgrade.exe" .\src\winget\upgrade\main.go

echo.
echo ^> Compiling slack-status (Go)...
go build -o ".\build\bin\slack-status.exe" .\src\slack\status\main.go

echo.
echo ^> Compiling slack-startup (Go)...
go build -o ".\build\bin\slack-startup.exe" .\src\slack\startup\main.go

echo.
echo ^> Compiling fscs (Rust)...
rustc -C strip=symbols -Clink-arg=/DEBUG:NONE -Clink-arg=/PDB:NONE .\src\functions\fscs.rs -o .\build\bin\fscs.exe

echo.
pause
