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

  Invoke-Expression (&shaka pwsh | Out-String)

  mise env --dotenv | ForEach-Object {
    if ($_ -match "^(.*?)=(.*)$") {
        $key = $matches[1]
        $val = $matches[2]
        ${env:$key} = $val
    }
  }

  Invoke-Expression (&zoxide init powershell | Out-String)
  Invoke-Expression (&starship init powershell)
}
