Remove-Item Alias:ni -Force -ErrorAction Ignore
Invoke-Expression (&starship init powershell)
Invoke-Expression (& { (zoxide init powershell | Out-String) })