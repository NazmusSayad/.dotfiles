@echo off

echo ^> Setting up git config...
call bash "./src/shell/git-config.sh"
git config --global credential.helper manager
git config --global core.sshCommand "C:/Windows/System32/OpenSSH/ssh.exe"

echo.
echo ^> Symlinking...
call go run ./src/scripts/symlink-init/main.go

echo.
echo ^> Installing tasks...
call go run ./src/install-windows-tasks/main.go

echo.
echo ^> Installing start menu entries...
call go run ./src/install-start-menu/main.go

echo.
echo Done!
