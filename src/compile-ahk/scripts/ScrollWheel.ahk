#UseHook
#NoTrayIcon
ProcessSetPriority "Realtime"
A_MaxHotkeysPerInterval := 9999

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