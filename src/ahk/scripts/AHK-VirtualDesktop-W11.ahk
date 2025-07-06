DLLPath := A_ScriptDir . "\..\src\ahk\bin\VirtualDesktopAccessor.dll"
GetDesktopCount := DllCall.Bind(DLLPath "\GetDesktopCount", "cdecl int")
GetCurrentDesktop := DllCall.Bind(DLLPath "\GetCurrentDesktopNumber", "cdecl int")

CapsLock::
{
  current := GetCurrentDesktop()
  total := GetDesktopCount()

  if (current >= total - 1) {
    Loop total - 1 {
      Send("^#{Left}")
      Sleep 50
    }
  } else {
    Send("^#{Right}")
  }
}
return