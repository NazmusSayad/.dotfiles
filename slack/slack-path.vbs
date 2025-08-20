Set objShell = CreateObject("WScript.Shell")
Set objNetwork = CreateObject("WScript.Network")

slackPath = "C:\Users\" & objNetwork.UserName & "\AppData\Local\slack"
slackBaseExePath = slackPath & "\slack.exe"


slackVersion = "4.45.69"
slackExePath = slackPath & "\app-" & slackVersion & "\slack.exe"

MsgBox "Slack Path: " & slackPath & vbCrLf & "Slack Version: " & slackVersion & vbCrLf & "Slack Exe: " & slackExePath

Set objShell = Nothing
