@echo off
setlocal

echo Linking current directory to user profile .dotfiles...

set "CURRENT_DIR=%~dp0"
set "CURRENT_DIR=%CURRENT_DIR:~0,-1%"
set "DOTFILES_DIR=%USERPROFILE%\.dotfiles"

if exist "%DOTFILES_DIR%" (
    echo Removing existing .dotfiles directory/link...
    rmdir /S /Q "%DOTFILES_DIR%" 2>nul
    if exist "%DOTFILES_DIR%" (
        del /F /Q "%DOTFILES_DIR%" 2>nul
    )
)

echo Creating symbolic link: %DOTFILES_DIR% == %CURRENT_DIR%
mklink /D "%DOTFILES_DIR%" "%CURRENT_DIR%"

if %errorLevel% EQU 0 (
    echo OK: .dotfiles linked successfully
) else (
    echo FAIL: Could not create symbolic link. Administrator privileges may be required.
    echo Press any key to exit...
    pause >nul
    exit /b 1
)

echo Press any key to continue...
pause >nul
