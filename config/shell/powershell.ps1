#!/usr/bin/env pwsh

Invoke-Expression (& dotsh pwsh @(mise env -D) | Out-String)

if ($PSVersionTable.PSEdition -eq 'Core') {
  Invoke-Expression (&pwshac cd | Out-String)
  Invoke-Expression (&shaka pwsh | Out-String)
  Invoke-Expression (&zoxide init powershell | Out-String)
  Invoke-Expression (&starship init powershell)
  
  zoxide add $PWD
}
