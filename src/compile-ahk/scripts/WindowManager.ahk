#UseHook
#NoTrayIcon
ProcessSetPriority "High"

RAlt & F19::AltTab
LAlt & F19::AltTab

F19:: {
  Send("{Alt down}")
  Send("{Tab}")
}

+F19:: {
  Send("{Shift down}")
  Send("{Alt down}")
  Send("{Tab}")
}

F19 up:: {
  Send("{Alt up}")
  Send("{Shift up}")
}

+F19 up:: {
  Send("{Alt up}")
  Send("{Shift up}")
}
