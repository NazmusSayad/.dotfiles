@echo off
setlocal

echo "> Installing shell providers..."
pacman -S zsh fish bash

echo "> Installing shell helpers..."
pacman -S mingw-w64-x86_64-starship mingw-w64-clang-x86_64-fastfetch

echo "> Installing helpers..."
pacman -S mingw-w64-x86_64-ffmpeg

echo.
echo "> Press any key to continue..."
pause >nul
