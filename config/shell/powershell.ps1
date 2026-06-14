if ($PSVersionTable.PSEdition -eq 'Core') {
  Invoke-Expression (& dotsh pwsh @(mise env -D) | Out-String)

  if (Test-Path ~/.env) { Invoke-Expression (& dotsh pwsh (Get-Content ~/.env -Raw) | Out-String) }
  if (Test-Path ~/.path) { $env:PATH += [IO.Path]::PathSeparator + ((Get-Content ~/.path) -join [IO.Path]::PathSeparator) }
  Invoke-Expression "$(direnv hook pwsh)"

  Invoke-Expression (&pwshac cd | Out-String)
  Invoke-Expression (&shaka pwsh | Out-String)
  Invoke-Expression (&zoxide init powershell | Out-String)
  Invoke-Expression (&starship init powershell)
}
