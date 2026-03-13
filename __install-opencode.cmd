@echo off
setlocal

if exist "%USERPROFILE%\.local\share\opencode\auth.json" (
    copy "%USERPROFILE%\.local\share\opencode\auth.json" "%TEMP%\auth.json.bak"
)

if exist "%USERPROFILE%\AppData\Roaming\ai.opencode.desktop\default.dat" (
    copy "%USERPROFILE%\AppData\Roaming\ai.opencode.desktop\default.dat" "%TEMP%\default.dat.bak"
)

rm -rf ~/.cache/opencode
rm -rf ~/.config/opencode
rm -rf ~/.local/share/opencode
rm -rf ~/.local/state/opencode

rm -rf ~/.bun/install/cache

rm -rf ~/AppData/Roaming/OpenCode
rm -rf ~/AppData/Local/ai.opencode.desktop
rm -rf ~/AppData/Roaming/ai.opencode.desktop

if exist "%TEMP%\auth.json.bak" (
    mkdir "%USERPROFILE%\.local\share\opencode" 2>nul
    move "%TEMP%\auth.json.bak" "%USERPROFILE%\.local\share\opencode\auth.json"
)

if exist "%TEMP%\default.dat.bak" (
    mkdir "%USERPROFILE%\AppData\Roaming\ai.opencode.desktop" 2>nul
    move "%TEMP%\default.dat.bak" "%USERPROFILE%\AppData\Roaming\ai.opencode.desktop\default.dat"
)

call sudo symlink-setup.exe

echo Done!
pause
