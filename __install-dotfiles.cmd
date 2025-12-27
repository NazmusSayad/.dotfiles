echo off
setlocal

set "CURRENT_DIR=%~dp0"
set "CURRENT_DIR=%CURRENT_DIR:~0,-1%"

set "DOTFILES_DIR=%CURRENT_DIR%"
set "DOTFILES_DIR_BIN=%DOTFILES_DIR%\.build\bin"

echo Setting DOTFILES_DIR...
powershell -NoProfile -Command "[Environment]::SetEnvironmentVariable('DOTFILES_DIR','%DOTFILES_DIR%','User')"

echo.
echo Checking if DOTFILES_DIR_BIN is in PATH...
for /f "usebackq delims=" %%P in (`
  powershell -NoProfile -Command "$p=[Environment]::GetEnvironmentVariable('Path','User') -split ';'; if($p -contains '%DOTFILES_DIR_BIN%'){'IN_PATH'} else {'NOT_IN_PATH'}"
`) do set "DOTFILES_DIR_BIN_IN_PATH=%%P"

echo DOTFILES_DIR_BIN_IN_PATH=%DOTFILES_DIR_BIN_IN_PATH%

if /I "%DOTFILES_DIR_BIN_IN_PATH%"=="NOT_IN_PATH" (
  echo DOTFILES_DIR_BIN is not in PATH.
  echo.
  echo Would set PATH to:
  powershell -NoProfile -Command "$p=[Environment]::GetEnvironmentVariable('Path','User'); if(-not ($p -split ';' | Where-Object { $_ -eq '%DOTFILES_DIR_BIN%' })) { $new = $p + ';%DOTFILES_DIR_BIN%'; Write-Output $new } else { Write-Output $p }"
) else (
  echo DOTFILES_DIR_BIN is already in PATH.
)
