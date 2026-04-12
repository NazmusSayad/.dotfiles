if ($PSVersionTable.PSEdition -eq 'Core') {
  $KEEP = "cd"

  $first = @(Get-Alias | Select-Object -ExpandProperty Name)
  foreach ($n in $first) {
    if ($KEEP -notcontains $n) {
      Remove-Item -LiteralPath "Alias:$n" -Force -ErrorAction SilentlyContinue
    }
  }
  $second = @(Get-Alias | Select-Object -ExpandProperty Name)
  foreach ($n in $second) {
    if ($KEEP -notcontains $n) {
      Remove-Item -LiteralPath "Alias:$n" -Force -ErrorAction SilentlyContinue
    }
  }

  Invoke-Expression (&shell-alias pwsh | Out-String)

  Invoke-Expression (&starship init powershell)
  Invoke-Expression (&zoxide init powershell | Out-String)
}
