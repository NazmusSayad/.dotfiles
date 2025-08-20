Set WshShell = CreateObject("WScript.Shell")
Set FSO = CreateObject("Scripting.FileSystemObject")

currentTime = Time()
currentHour = Hour(currentTime)
currentDay = WeekDay(Date())

' Check if current time is between 6:00 AM and 8:00 PM
isWithinTimeRange = (currentHour >= 6) And (currentHour < 20)

' Check if it's not Friday (6) or Saturday (7)
isNotWeekend = (currentDay <> 6) And (currentDay <> 7)

' Check if Slack is already running
Set objWMIService = GetObject("winmgmts:")
Set colProcesses = objWMIService.ExecQuery("Select * From Win32_Process Where Name = 'slack.exe'")
isSlackRunning = (colProcesses.Count > 0)

If isWithinTimeRange And isNotWeekend And Not isSlackRunning Then
  WshShell.Run """C:\Users\Sayad\AppData\Local\slack\app-4.45.69\slack.exe"" --startup --silent", 0, False
End If

Set WshShell = Nothing
Set FSO = Nothing
Set objWMIService = Nothing