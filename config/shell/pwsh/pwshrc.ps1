if ([Environment]::CommandLine -match '-NonI') { return }

pwshac cd | Out-String -Width ([int]::MaxValue) | Invoke-Expression
mise activate pwsh | Out-String -Width ([int]::MaxValue) | Invoke-Expression

$dotsysDir = if ($env:DOTSYS_DIR) { $env:DOTSYS_DIR } else { "$HOME/.dotsys" }
dotsh init pwsh --home $dotsysDir | Out-String -Width ([int]::MaxValue) | Invoke-Expression

direnv hook pwsh | Out-String -Width ([int]::MaxValue) | Invoke-Expression

shaka pwsh | Out-String -Width ([int]::MaxValue) | Invoke-Expression
starship init powershell | Out-String -Width ([int]::MaxValue) | Invoke-Expression
zoxide init powershell | Out-String -Width ([int]::MaxValue) | Invoke-Expression
zoxide add "$PWD" 2>$null
