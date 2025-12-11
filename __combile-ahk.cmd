@echo off

echo ^> Deleting old compiled files...
del .\build\ahk /s /q

echo.
echo ^> Compiling ahk scripts...
call .\build\bin\ahk-compile.exe

echo.
pause
