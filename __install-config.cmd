@echo off
setlocal

echo Setting up git config...
git config --global user.name "Nazmus Sayad"
git config --global user.email "87106526+NazmusSayad@users.noreply.github.com"
git config --global init.defaultBranch main
git config --global --add safe.directory "*"
git config --global --add --bool push.autoSetupRemote true
git config --global core.eol lf
git config --global core.autocrlf false
git config --global core.pager cat
git config --global core.ignorecase false
git config --global core.editor "code --wait"

echo.
echo Symlinking...
call symlink-setup.exe

echo.
echo Installing tasks...
call go run ./src/install-windows-tasks/main.go

echo.
echo Installing start menu entries...
call go run ./src/install-start-menu/main.go

echo.
echo Installing Mise...
mise install

echo.
echo Installing MSYS2 shells...
pacman -S --noconfirm bash fish zsh

echo.
echo Installing MSYS2 tools...
pacman -S --noconfirm nnn ncdu mingw-w64-x86_64-gdu

echo.
echo Setting up GOPATH...
call powershell.exe -NoProfile -Command "[Environment]::SetEnvironmentVariable('GOPATH', (go env GOPATH), 'User')"

echo.
echo Setting up GOROOT...
call powershell.exe -NoProfile -Command "[Environment]::SetEnvironmentVariable('GOROOT', (go env GOROOT), 'User')"

echo.
pause
