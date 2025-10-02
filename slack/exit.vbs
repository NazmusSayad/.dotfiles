Set WshShell = CreateObject("WScript.Shell")
Set FSO = CreateObject("Scripting.FileSystemObject")

currentTime = Time()
currentHour = Hour(currentTime)
currentDay = WeekDay(Date())

' Check if current time is outside 6:00 AM and 8:00 PM
isOutsideTimeRange = (currentHour < 6) Or (currentHour >= 20)

' Check if it's Friday (6) or Saturday (7)
isWeekend = (currentDay = 6) Or (currentDay = 7)

' Check if Slack is running
Set objWMIService = GetObject("winmgmts:")
Set colProcesses = objWMIService.ExecQuery("Select * From Win32_Process Where Name = 'slack.exe'")
isSlackRunning = (colProcesses.Count > 0)

If (isOutsideTimeRange Or isWeekend) And isSlackRunning Then
  WshShell.Run "taskkill /f /im slack.exe", 0, True
End If

Set WshShell = Nothing
Set FSO = Nothing
Set objWMIService = Nothing
