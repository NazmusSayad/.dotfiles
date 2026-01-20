#UseHook
#NoTrayIcon
ProcessSetPriority "High"

^#!F11:: {
  hwnd := WinGetID("A")
  if !hwnd
    return

  title := WinGetTitle(hwnd)
  ex    := WinGetExStyle(hwnd)

  WS_EX_TOOLWINDOW := 0x80
  WS_EX_APPWINDOW  := 0x40000

  if (ex & WS_EX_TOOLWINDOW) {
    WinSetExStyle "-" WS_EX_TOOLWINDOW, hwnd
    WinSetExStyle "+" WS_EX_APPWINDOW, hwnd
    TrayTip "Window Restored", title, 1500
  } else {
    WinSetExStyle "+" WS_EX_TOOLWINDOW, hwnd
    WinSetExStyle "-" WS_EX_APPWINDOW, hwnd
    TrayTip "Window Hidden", title, 1500
  }
}
