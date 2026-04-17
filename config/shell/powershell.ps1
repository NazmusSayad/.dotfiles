Invoke-Expression (&mise env -s pwsh | Out-String)
Invoke-Expression (&zoxide init powershell | Out-String)

if ($PSVersionTable.PSEdition -eq 'Core') {
  Invoke-Expression (&pwshac cd | Out-String)
  Invoke-Expression (&shaka pwsh | Out-String)

  Invoke-Expression (&starship init powershell)
}
