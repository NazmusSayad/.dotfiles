@echo off
powershell -Command "if (-not([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] 'Administrator')) { exit 1 }"
if %errorLevel% NEQ 0 (
    echo This script requires administrator privileges.
    echo Press any key to exit...
    pause >nul
    exit /b
)

setlocal

set __dirname=%~dp0
set __dirname=%__dirname:~0,-1%
echo CWD: %__dirname%
echo.

if not exist "%USERPROFILE%\.config" mkdir "%USERPROFILE%\.config"
echo User config directory created
echo.

if exist "%USERPROFILE%\.config\fastfetch" rmdir /S /Q "%USERPROFILE%\.config\fastfetch"
mklink /D "%USERPROFILE%\.config\fastfetch" "%__dirname%\config\fastfetch"
echo Fastfetch config linked
echo.

if exist "%USERPROFILE%\.config\starship.toml" del "%USERPROFILE%\.config\starship.toml"
mklink "%USERPROFILE%\.config\starship.toml" "%__dirname%\config\starship.toml"
echo Starship config linked
echo.

if not exist "%USERPROFILE%\.config\fish" mkdir "%USERPROFILE%\.config\fish"
echo Fish config directory created
echo.

if exist "%USERPROFILE%\.config\fish\config.fish" del "%USERPROFILE%\.config\fish\config.fish"
mklink "%USERPROFILE%\.config\fish\config.fish" "%__dirname%\config\fish-config.fish"
echo Fish config linked
echo.

if exist "%USERPROFILE%\.zshrc" del "%USERPROFILE%\.zshrc"
mklink "%USERPROFILE%\.zshrc" "%__dirname%\config\zsh-config.sh"
echo Zsh config linked
echo.

if exist "%USERPROFILE%\.bashrc" del "%USERPROFILE%\.bashrc"
mklink "%USERPROFILE%\.bashrc" "%__dirname%\config\bash-config.sh"
echo Bash config linked
echo.

git config --global user.name "Nazmus Sayad"
git config --global user.email "87106526+NazmusSayad@users.noreply.github.com"

git config --global init.defaultBranch main

git config --global core.eol lf
git config --global core.autocrlf false

git config --global core.pager cat
git config --global core.ignorecase false

git config --global --add safe.directory "*"
git config --global --add --bool push.autoSetupRemote true

echo Git config added

echo.
pause
