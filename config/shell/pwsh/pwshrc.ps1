if ([Environment]::CommandLine -match '-NonI') { return }

pwshac cd | Out-String -Width ([int]::MaxValue) | Invoke-Expression
mise activate pwsh | Out-String -Width ([int]::MaxValue) | Invoke-Expression

dotsh init pwsh --home "$HOME/.dotsys" | Out-String -Width ([int]::MaxValue) | Invoke-Expression

direnv hook pwsh | Out-String -Width ([int]::MaxValue) | Invoke-Expression

shaka pwsh | Out-String -Width ([int]::MaxValue) | Invoke-Expression
starship init powershell | Out-String -Width ([int]::MaxValue) | Invoke-Expression
zoxide init powershell | Out-String -Width ([int]::MaxValue) | Invoke-Expression
zoxide add "$PWD" 2>$null
