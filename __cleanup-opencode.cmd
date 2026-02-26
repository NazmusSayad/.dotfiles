@echo off
setlocal

rm -rf ~/.config/opencode
rm -rf ~/.local/share/opencode

rm -rf ~/.bun/install/cache

rm -rf ~/AppData/Roaming/OpenCode
rm -rf ~/AppData/Local/ai.opencode.desktop
rm -rf ~/AppData/Roaming/ai.opencode.desktop

echo Done!
pause
