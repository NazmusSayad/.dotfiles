@echo off
setlocal

set "CURRENT_DIR=%~dp0"
set "CURRENT_DIR=%CURRENT_DIR:~0,-1%"
set "DOTFILES_DIR=%USERPROFILE%\.dotfiles"

:: Remove existing .dotfiles directory/link
if exist "%DOTFILES_DIR%" (
    echo Removing existing .dotfiles directory/link...
    rmdir /S /Q "%DOTFILES_DIR%" 2>nul
    if exist "%DOTFILES_DIR%" (
        del /F /Q "%DOTFILES_DIR%" 2>nul
    )
)

setx DOTFILES_DIR "%CURRENT_DIR%"
setx PATH "%CURRENT_DIR%\.build\bin;%PATH%"

echo.
echo SUCCESS: Dotfiles installation completed.
pause