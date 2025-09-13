Set WshShell = CreateObject("WScript.Shell")
WshShell.Run """cmd.exe"" /c", 0, False
Set WshShell = Nothing

Set UAC = CreateObject("Shell.Application")
UAC.ShellExecute "F:\___SYSTEM___\dotfiles\___AHK-Macro.exe", "", "", "runas", 0
Set UAC = Nothing

Set UAC = CreateObject("Shell.Application")
UAC.ShellExecute "F:\___SYSTEM___\dotfiles\___AHK-VirtualDesktop-W11.exe", "", "", "runas", 0
Set UAC = Nothing

Set UAC = CreateObject("Shell.Application")
UAC.ShellExecute "C:\Program Files\ShareX\ShareX.exe", "", "", "runas", 0
Set UAC = Nothing

Set WshShell = CreateObject("WScript.Shell")
WshShell.Run """gpg"" --list-keys", 0, False
Set WshShell = Nothing