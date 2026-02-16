#UseHook
#NoTrayIcon
ProcessSetPriority "High"

RAlt & F23::AltTab
LAlt & F23::AltTab

F23:: {
  Send("{Alt down}")
  Send("{Tab}")
}

F23 up:: {
  Send("{Alt up}")
}
