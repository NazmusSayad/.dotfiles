#NoTrayIcon

Initilized := false
F23::HandleAutoDesktopSwitch

HandleAutoDesktopSwitch() { 
  global Initilized, CurrentDesktop
  If (!Initilized) {
    Initilized := true
    Return Send("^#{Right}")
  }

  mapDesktopsFromRegistry()
  If (CurrentDesktop >= DesktopCount) {
    CurrentDesktop := 1
  } Else {
    CurrentDesktop++
  }

  switchDesktopByNumber(CurrentDesktop)
}

mapDesktopsFromRegistry() {
  global DesktopCount, CurrentDesktop
  DesktopCount := 1

  IdLength := 32
  SessionId := getSessionId()
  If (SessionId) {
    CurrentDesktopId := RegRead("HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\SessionInfo\" SessionId "\VirtualDesktops", "CurrentVirtualDesktop")
    If (CurrentDesktopId) {
      IdLength := StrLen(CurrentDesktopId)
    }
  }

  DesktopList := RegRead("HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\VirtualDesktops", "VirtualDesktopIDs")
  If (DesktopList) {
    DesktopListLength := StrLen(DesktopList)
    DesktopCount := DesktopListLength / IdLength
  }

  i := 0
  While (CurrentDesktopId And i < DesktopCount) {
    StartPos := (i * IdLength) + 1
    DesktopIter := SubStr(DesktopList, StartPos, IdLength)
    If (DesktopIter = CurrentDesktopId) {
      CurrentDesktop := i + 1
      Break
    }
    i++
  }
}

getSessionId() {
  ProcessId := DllCall("GetCurrentProcessId", "UInt")
  If (ProcessId = 0) {
    OutputDebug("Error getting current process id.")
    Return 0
  }
  SessionId := 0
  Result := DllCall("ProcessIdToSessionId", "UInt", ProcessId, "UInt*", &SessionId)
  If (Result = 0) {
    OutputDebug("Error getting session id.")
    Return 0
  }

  Return SessionId
}

switchDesktopByNumber(TargetDesktop) {
  global CurrentDesktop
  mapDesktopsFromRegistry()

  If (TargetDesktop > DesktopCount Or TargetDesktop < 1) {
    Return OutputDebug("[invalid] target: " TargetDesktop " current: " CurrentDesktop)
  }

  SleepTime := 100
  TargetDiff := TargetDesktop - CurrentDesktop
  If (TargetDiff != 0) {
    SleepTime := Round(100 / Max(TargetDiff, TargetDiff * -1, 1))
  } 

  While (CurrentDesktop < TargetDesktop) {
    Send "^#{Right}"
    Sleep(SleepTime)
    CurrentDesktop++
  }

  While (CurrentDesktop > TargetDesktop) {
    Send "^#{Left}"
    Sleep(SleepTime)
    CurrentDesktop--
  }
}
