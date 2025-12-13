@echo off
setlocal

echo ^> Setting up git config...
git config --global user.name "Nazmus Sayad"
git config --global user.email "87106526+NazmusSayad@users.noreply.github.com"
git config --global init.defaultBranch main
git config --global core.eol lf
git config --global core.autocrlf false
git config --global core.pager cat
git config --global core.ignorecase false
git config --global --add safe.directory "*"
git config --global --add --bool push.autoSetupRemote true

echo ^> pnpm config settings...
call pnpm config set ci true
call pnpm config set allow-scripts true
call pnpm config set shamefully-hoist true
call pnpm config set auto-install-peers true

echo ^> Symlinking
call symlink-setup.exe

echo.
pause