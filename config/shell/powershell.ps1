Remove-Item Alias:r -Force -ErrorAction Ignore
Remove-Item Alias:ni -Force -ErrorAction Ignore
Remove-Item Alias:gc -Force -ErrorAction Ignore

Invoke-Expression (&shell-alias pwsh | Out-String)

if ($PSVersionTable.PSEdition -eq 'Core') {
  Invoke-Expression (&starship init powershell)
}
