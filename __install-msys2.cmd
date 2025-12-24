@echo off
setlocal

echo Installing shells...
pacman -S --noconfirm bash fish zsh

echo Installing tools...
pacman -S --noconfirm lua nnn ncdu mingw-w64-ucrt-x86_64-dust

echo.
echo Done. Press any key to exit...
pause>nul
