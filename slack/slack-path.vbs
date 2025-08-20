Set objShell = CreateObject("WScript.Shell")
Set objNetwork = CreateObject("WScript.Network")

slackPath = "C:\Users\" & objNetwork.UserName & "\AppData\Local\slack"
slackBaseExePath = slackPath & "\slack.exe"

Set objExec = objShell.Exec("powershell -Command ""(Get-Item '" & slackBaseExePath & "').VersionInfo.FileVersion""")
slackVersion = Trim(objExec.StdOut.ReadAll)
slackExePath = slackPath & "\app-" & slackVersion & "\slack.exe"

MsgBox "Slack Path: " & slackPath & vbCrLf & "Slack Version: " & slackVersion & vbCrLf & "Slack Exe: " & slackExePath
Set objShell = Nothing
