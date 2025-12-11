@echo off
setlocal

echo "> Installing shells..."
pacman -S zsh fish bash

echo "> Installing fastfetch..."
pacman -S mingw-w64-clang-x86_64-fastfetch

echo "> Installing ffmpeg..."
pacman -S mingw-w64-x86_64-ffmpeg

echo "> Fish shell config..."
set fish_color_command magenta

echo.
echo "> Press any key to continue..."
pause >nul
