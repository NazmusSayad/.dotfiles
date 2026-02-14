#UseHook
#NoTrayIcon
ProcessSetPriority "Realtime"
A_MaxHotkeysPerInterval := 9999

#Space::+!F

#PrintScreen::#^+PrintScreen

^WheelUp:: {
  Critical "On"
  KeyWait "Ctrl"
  Send "{WheelUp}"
}

^WheelDown:: {
  Critical "On"
  KeyWait "Ctrl"
  Send "{WheelDown}"
}
 
#`:: {
  RunWait 'schtasks /create /tn "WT-Desktop-User" /tr "wt.exe -d \"' A_Desktop '\"" /sc once /st 00:00 /sd 01/01/2000 /f', , "Hide"
  Run 'schtasks /run /tn "WT-Desktop-User"', , "Hide"
}

^#F11:: {
  SoundBeep 1000, 300
  WinSetAlwaysOnTop(-1, "A")
}

^#!F11:: {
  hwnd := WinGetID("A")
  if !hwnd {
    return
  }

  title := WinGetTitle(hwnd)
  ex := WinGetExStyle(hwnd)

  WS_EX_TOOLWINDOW := 0x80
  WS_EX_APPWINDOW := 0x40000

  SoundBeep 1000, 300
  if (ex & WS_EX_TOOLWINDOW) {
    WinSetExStyle "-" WS_EX_TOOLWINDOW, hwnd
    WinSetExStyle "+" WS_EX_APPWINDOW, hwnd
  } else {
    WinSetExStyle "+" WS_EX_TOOLWINDOW, hwnd
    WinSetExStyle "-" WS_EX_APPWINDOW, hwnd
  }
}

::@me::me@sayad.dev
::@fake::fake@sayad.dev
::@env::sayadenv@gmail.com
::@mail::247sayad@gmail.com
