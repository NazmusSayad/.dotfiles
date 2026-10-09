if ([Environment]::CommandLine -match '-NonI') { return }

pwshac | Out-String -Width ([int]::MaxValue) | Invoke-Expression
dotsh pwsh (mise env --dotenv) | Out-String -Width ([int]::MaxValue) | Invoke-Expression

$dotfilesDir = if ($env:DOTFILES_DIR) { $env:DOTFILES_DIR } else { "$HOME/.dotfiles" }
dotsh init pwsh --home $dotfilesDir | Out-String -Width ([int]::MaxValue) | Invoke-Expression

direnv hook pwsh | Out-String -Width ([int]::MaxValue) | Invoke-Expression

shaka pwsh | Out-String -Width ([int]::MaxValue) | Invoke-Expression
zoxide init powershell | Out-String -Width ([int]::MaxValue) | Invoke-Expression
zoxide add "$PWD" 2>$null
starship init powershell | Out-String -Width ([int]::MaxValue) | Invoke-Expression
