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


::@me::me@sayad.dev
::@fake::fake@sayad.dev
::@env::sayadenv@gmail.com
::@mail::247sayad@gmail.com
