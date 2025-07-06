@echo off
setlocal

REM Initialize
set __dirname=%~dp0
set __dirname=%__dirname:~0,-1%
echo CWD: %__dirname%
echo.

REM Setting up fish config
if not exist "%USERPROFILE%\.config\fish" mkdir "%USERPROFILE%\.config\fish"
echo source "%__dirname%/config/fish-config/__init__.fish" > "%USERPROFILE%\.config\fish\config.fish"
echo Fish config linked
echo.

REM Git configuration
git config --global user.name "Nazmus Sayad"
git config --global user.email "87106526+NazmusSayad@users.noreply.github.com"

git config --global core.eol lf
git config --global core.autocrlf false

git config --global core.pager cat
git config --global core.ignorecase false

git config --global init.defaultBranch main
git config --global core.excludesfile ~/default.gitignore

git config --global --add safe.directory "*"
git config --global --add --bool push.autoSetupRemote true

echo Git config added

REM Setting up npm config
npm config set ignore-scripts true
echo Npm config added

echo.
echo Press any key to continue...
pause