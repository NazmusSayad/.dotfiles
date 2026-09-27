#UseHook
#NoTrayIcon
ProcessSetPriority "High"

RAlt & F20::AltTab
LAlt & F20::AltTab

F20:: {
  Send("{Alt down}")
  Send("{Tab}")
}

+F20:: {
  Send("{Shift down}")
  Send("{Alt down}")
  Send("{Tab}")
}

F20 up:: {
  Send("{Alt up}")
  Send("{Shift up}")
}

+F20 up:: {
  Send("{Alt up}")
  Send("{Shift up}")
}
