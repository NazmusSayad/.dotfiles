Get-Alias | ForEach-Object {
  Remove-Alias $_.Name -Force -ErrorAction Ignore
}

Invoke-Expression (&shell-alias pwsh | Out-String)

if ($PSVersionTable.PSEdition -eq 'Core') {
  Invoke-Expression (&starship init powershell)
}
