#UseHook
#NoTrayIcon

::@me::me@sayad.dev
::@fake::fake@sayad.dev
::@env::sayadenv@gmail.com
::@mail::247sayad@gmail.com

#Space::+!F
#PrintScreen::#^+PrintScreen

#`:: {
  Run 'gsudo --integrity Medium powershell -Command "Start-Process wt -WorkingDirectory $env:USERPROFILE\\Desktop"'
    , , "Hide"
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
