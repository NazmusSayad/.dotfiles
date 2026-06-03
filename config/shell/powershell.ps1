if ($PSVersionTable.PSEdition -eq 'Core') {
  Invoke-Expression (& dotsh pwsh @(mise env -D) | Out-String)
  Invoke-Expression "$(direnv hook pwsh)"

  Invoke-Expression (&pwshac cd | Out-String)
  Invoke-Expression (&shaka pwsh | Out-String)
  Invoke-Expression (&zoxide init powershell | Out-String)
  Invoke-Expression (&starship init powershell)
}
