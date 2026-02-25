#NoTrayIcon

DLLPath := A_ScriptDir . "\..\lib\VirtualDesktopAccessor.dll"
GetDesktopCount := DllCall.Bind(DLLPath "\GetDesktopCount", "cdecl int")
GetCurrentDesktop := DllCall.Bind(DLLPath "\GetCurrentDesktopNumber", "cdecl int")

^F23:: {
  current := GetCurrentDesktop()
  total := GetDesktopCount()

  if (current >= total - 1) {
    loop total - 1 {
      Send("^#{Left}")
    }
  } else {
    Send("^#{Right}")
  }
}