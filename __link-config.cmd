@echo off
setlocal

set __dirname=%~dp0
set __dirname=%__dirname:~0,-1%
echo CWD: %__dirname%
echo.

if not exist "%USERPROFILE%\.config" mkdir "%USERPROFILE%\.config"
echo "> User config directory created..."
echo.

if exist "%USERPROFILE%\.config\fastfetch" rmdir /S /Q "%USERPROFILE%\.config\fastfetch"
mklink /D "%USERPROFILE%\.config\fastfetch" "%__dirname%\config\fastfetch"
echo "> Fastfetch config linked..."
echo.

if exist "%USERPROFILE%\.config\starship.toml" del "%USERPROFILE%\.config\starship.toml"
mklink "%USERPROFILE%\.config\starship.toml" "%__dirname%\config\starship.toml"
echo "> Starship config linked..."
echo.

if not exist "%USERPROFILE%\.config\fish" mkdir "%USERPROFILE%\.config\fish"
echo "> Fish config directory created..."
echo.

if exist "%USERPROFILE%\.config\fish\config.fish" del "%USERPROFILE%\.config\fish\config.fish"
mklink "%USERPROFILE%\.config\fish\config.fish" "%__dirname%\config\fish-config.fish"
echo "> Fish config linked..."
echo.

if exist "%USERPROFILE%\.zshrc" del "%USERPROFILE%\.zshrc"
mklink "%USERPROFILE%\.zshrc" "%__dirname%\config\zsh-config.sh"
echo "> Zsh config linked..."
echo.

if exist "%USERPROFILE%\.bashrc" del "%USERPROFILE%\.bashrc"
mklink "%USERPROFILE%\.bashrc" "%__dirname%\config\bash-config.sh"
echo "> Bash config linked..."
echo.

echo.
echo "> Press any key to continue..."
pause >nul
