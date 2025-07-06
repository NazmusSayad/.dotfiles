Set WshShell = CreateObject("WScript.Shell")
WshShell.Run """cmd.exe"" /c", 0, False
Set WshShell = Nothing

Set WshShell = CreateObject("WScript.Shell")
WshShell.Run """E:\!#SYSTEM\.ahk-macro\bin\AHK-Macro.exe""", 0, False
Set WshShell = Nothing

Set WshShell = CreateObject("WScript.Shell")
WshShell.Run """E:\!#SYSTEM\.ahk-macro\bin\AHK-VirtualDesktop-W10.exe""", 0, False
Set WshShell = Nothing

Set WshShell = CreateObject("WScript.Shell")
WshShell.Run """C:\Program Files\ShareX\ShareX.exe""", 0, False
Set WshShell = Nothing

Set WshShell = CreateObject("WScript.Shell")
WshShell.Run """gpg"" --list-keys", 0, False
Set WshShell = Nothing
