@echo off
setlocal

echo Installing npm packages globally
volta install npm pnpm tsx uni-run code-info netserv

echo.
echo Press any key to continue...
pause >nul
