@echo off
setlocal

rm -rf ~/.cache/opencode
rm -rf ~/.config/opencode
rm -rf ~/.local/share/opencode
rm -rf ~/.local/state/opencode

rm -rf ~/.bun/install/cache

rm -rf ~/AppData/Roaming/OpenCode
rm -rf ~/AppData/Roaming/@opencode-ai
rm -rf ~/AppData/Local/@opencode-aidesktop-electron-updater
rm -rf ~/AppData/Local/ai.opencode.desktop
rm -rf ~/AppData/Local/ai.opencode.desktop.beta
rm -rf ~/AppData/Roaming/ai.opencode.desktop
rm -rf ~/AppData/Roaming/ai.opencode.desktop.beta

rm -rf ~/.t3
rm -rf ~/AppData/Roaming/t3code
rm -rf ~/AppData/Local/t3code-updater

call sudo symlink-init.exe
call proxy-pause echo
