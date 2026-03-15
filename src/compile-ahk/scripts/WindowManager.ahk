#UseHook
#NoTrayIcon
ProcessSetPriority "High"

RAlt & F19::AltTab
LAlt & F19::AltTab

F19:: {
  Send("{Alt down}")
  Send("{Tab}")
}

F19 up:: {
  Send("{Alt up}")
}
