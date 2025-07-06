@echo off
set AHK_EXEC=.\src\bin\AutoHotkey64.exe
set SCRIPTS_DIR=.\src\scripts

taskkill /IM AutoHotkey64.exe /F

for %%f in (%SCRIPTS_DIR%\*.ahk) do (
    echo Running %%f
    start "" run exec %AHK_EXEC% "%%f" -- --ext ahk
)
