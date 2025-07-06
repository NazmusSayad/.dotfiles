@echo off
setlocal enabledelayedexpansion

:: Enable ANSI color support
for /f "tokens=2 delims=[]" %%I in ('ver') do set winver=%%I
for /f "tokens=2,3,4 delims=. " %%I in ("%winver%") do (
    if %%I geq 10 (
        :: Windows 10+ supports ANSI escape sequences
        reg add HKEY_CURRENT_USER\Console /v VirtualTerminalLevel /t REG_DWORD /d 1 /f >nul 2>&1
    )
)

:: Set console size (width=80, height=45)
mode con cols=80 lines=45

:: Check if GPG is installed
where gpg >nul 2>&1
if !errorlevel! neq 0 (
    echo Error: GPG not installed
    exit /b 1
)

:: Check if Git is installed
where git >nul 2>&1
if !errorlevel! neq 0 (
    echo Error: Git not installed
    exit /b 1
)

:: Get Git configuration
for /f "tokens=*" %%i in ('git config --get user.name 2^>nul') do set "git_name=%%i"
for /f "tokens=*" %%i in ('git config --get user.email 2^>nul') do set "git_email=%%i"

if "%git_email%"=="" (
    echo Error: Git user.email or user.name not configured
    echo Please run: git config --global user.name "Your Name"
    echo Please run: git config --global user.email "your@email.com"
    exit /b 1
)

if "%git_name%"=="" (
    echo Error: Git user.email or user.name not configured
    echo Please run: git config --global user.name "Your Name"
    echo Please run: git config --global user.email "your@email.com"
    exit /b 1
)

echo Git user name  : %git_name%
echo Git user email : %git_email%

:: Check if GPG keys exist
gpg --list-secret-keys --keyid-format LONG 2>nul | findstr "sec" >nul
if !errorlevel! neq 0 (
    :: Create temporary batch file for GPG key generation
    set "batch_file=%temp%\gpg_batch_%random%.txt"
    (
        echo Key-Type: RSA
        echo Key-Length: 4096
        echo Key-Usage: sign
        echo Name-Real: %git_name%
        echo Name-Email: %git_email%
        echo Expire-Date: 0
        echo %%no-protection
        echo %%commit
    ) > "!batch_file!"

    gpg --batch --generate-key "!batch_file!"
    del "!batch_file!"

    if !errorlevel! neq 0 (
        exit /b 1
    )
)

:: Get GPG key ID
set "gpg_key_id="
for /f "tokens=*" %%i in ('gpg --list-secret-keys --keyid-format LONG 2^>nul') do (
    echo %%i | findstr "sec" >nul
    if !errorlevel! equ 0 (
        for /f "tokens=2 delims=/" %%j in ("%%i") do (
            for /f "tokens=1 delims= " %%k in ("%%j") do (
                set "gpg_key_id=%%k"
                goto :found_key
            )
        )
    )
)
:found_key

if "%gpg_key_id%"=="" (
    exit /b 1
)

:: Configure Git with GPG
git config --global user.signingkey %gpg_key_id%
git config --global commit.gpgsign true
for /f "tokens=*" %%i in ('where gpg') do git config --global gpg.program "%%i"

echo GPG key ID     : %gpg_key_id%

echo GPG key generated and configured for Git.

echo.
echo.
gpg --armor --export %gpg_key_id%
echo.
echo.

echo Press any key to exit...
pause >nul
