Set WshShell = CreateObject("WScript.Shell")

' Check if Slack is running
Set objWMIService = GetObject("winmgmts:")
Set colProcesses = objWMIService.ExecQuery("Select * From Win32_Process Where Name = 'slack.exe'")
isSlackRunning = colProcesses.Count > 0

If isSlackRunning Then
  WshShell.Run "taskkill /f /im slack.exe", 0, True
End If

Set WshShell = Nothing
Set objWMIService = Nothing
