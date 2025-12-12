#NoTrayIcon
ProcessSetPriority "Realtime"
A_MaxHotkeysPerInterval := 9999

#UseHook
#Space::+!F

#UseHook
#PrintScreen::#^+PrintScreen

#UseHook
^WheelUp:: {
  Critical "On"
  KeyWait "Ctrl"
  Send "{WheelUp}"
}

#UseHook
^WheelDown:: {
  Critical "On"
  KeyWait "Ctrl"
  Send "{WheelDown}"
}


::@me::me@sayad.dev
::@fake::fake@sayad.dev
::@mail::sayadenv@gmail.com
