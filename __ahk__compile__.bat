@echo off

call npm run ahk-compile
echo ahk-compile exit code: %errorLevel%

echo.
pause