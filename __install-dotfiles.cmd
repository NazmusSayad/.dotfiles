@echo off
setlocal

net session >nul 2>&1
if %errorLevel% NEQ 0 (
    echo FAIL: Administrator privileges required.
    echo Relaunching with elevated privileges...
    powershell -NoProfile -Command "Start-Process -FilePath '%~f0' -Verb RunAs"
    exit /b
)

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

:: Create symbolic link
mklink /D "%DOTFILES_DIR%" "%CURRENT_DIR%"
if %errorLevel% NEQ 0 (
    echo FAIL: Could not create symbolic link.
    echo Press any key to exit...
    pause >nul
    exit /b 1
)

:: Set up PATH environment variable
echo.
echo Setting up PATH environment variable...
set "DOTFILES_BIN=%DOTFILES_DIR%\build\bin"

echo %PATH% | find /i "%DOTFILES_BIN%" >nul
if %errorLevel% EQU 0 (
    echo OK: %DOTFILES_BIN% already in PATH.
) else (
    setx Path "%DOTFILES_BIN%;%PATH%" /M >nul

    if %errorLevel% EQU 0 (
        echo OK: PATH updated successfully.
    ) else (
        echo FAIL: Could not update PATH.
        echo Press any key to exit...
        pause >nul
        exit /b 1
    )
)

echo.
echo SUCCESS: Dotfiles installation completed.
echo Press any key to exit...
pause >nul