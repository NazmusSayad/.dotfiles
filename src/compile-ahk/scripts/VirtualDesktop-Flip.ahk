#NoTrayIcon

!F23:: {
  static goRight := true

  if (goRight) {
    Send("^#{Right}")
  } else {
    Send("^#{Left}")
  }

  goRight := !goRight
}