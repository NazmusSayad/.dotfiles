Set ObjWMIService = GetObject("winmgmts:")

currentTime = Time()
currentHour = Hour(currentTime)
currentDay = WeekDay(Date())

' Check if current time is between 6:00 AM and 8:00 PM
isWithinTimeRange = (currentHour >= 6) And (currentHour < 20)

' Check if it's not Friday (6) or Saturday (7)
isNotWeekend = (currentDay <> 6) And (currentDay <> 7)

' Check if Slack is already running
Set ColProcesses = ObjWMIService.ExecQuery("Select * From Win32_Process Where Name = 'slack.exe'")
isSlackRunning = (ColProcesses.Count > 0)

If isWithinTimeRange And isNotWeekend And Not isSlackRunning Then
  Set WshShell = CreateObject("WScript.Shell")
  Set PathObjShell = CreateObject("WScript.Shell")
  Set PathObjNetwork = CreateObject("WScript.Network")
  Set PathFSObject = CreateObject("Scripting.FileSystemObject")

  slackPath = "C:\Users\" & PathObjNetwork.UserName & "\AppData\Local\slack"
  slackBaseExePath = slackPath & "\slack.exe"
  tempFile = PathFSObject.GetSpecialFolder(2) & "\slack_ver.tmp"

  psCommand = "powershell -NoProfile -WindowStyle Hidden -Command ""(Get-Item '" & slackBaseExePath & "').VersionInfo.ProductVersion | Out-File -Encoding ASCII '" & tempFile & "'"""
  PathObjShell.Run psCommand, 0, True

  Set file = PathFSObject.OpenTextFile(tempFile, 1)
  slackVersion = Trim(file.ReadAll)
  slackVersion = Replace(slackVersion, vbCrLf, "")
  slackVersion = Replace(slackVersion, vbCr, "")
  slackVersion = Replace(slackVersion, vbLf, "")
  slackVersion = Trim(slackVersion)

  file.Close
  PathFSObject.DeleteFile tempFile

  slackExePath = slackPath & "\app-" & slackVersion & "\slack.exe"
  WshShell.Run """" & slackExePath & """ --startup", 0, False

  Set file = Nothing
  Set WshShell = Nothing
  Set PathObjShell = Nothing
  Set PathObjNetwork = Nothing
  Set PathFSObject = Nothing
End If

Set ObjWMIService = Nothing
Set ColProcesses = Nothing