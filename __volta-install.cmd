@echo off
setlocal

echo Installing npm packages globally
volta install node@latest npm@latest pnpm@latest tsx@latest uni-run@latest code-info@latest netserv@latest

echo.
echo Press any key to continue...
pause >nul
