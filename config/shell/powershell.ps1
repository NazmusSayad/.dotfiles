if ($PSVersionTable.PSEdition -eq 'Core') {
  Invoke-Expression (&pwshac cd | Out-String)
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
