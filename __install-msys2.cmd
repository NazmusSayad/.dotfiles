@echo off
setlocal

echo ^> Installing shells...
pacman -S --noconfirm bash fish

echo ^> Installing fastfetch...
pacman -S --noconfirm mingw-w64-clang-x86_64-fastfetch

echo.
pause
