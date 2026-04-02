Remove-Item Alias:ni -Force -ErrorAction Ignore

if ($PSVersionTable.PSEdition -eq 'Core') {
  Invoke-Expression (&starship init powershell)
}