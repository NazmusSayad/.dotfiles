if ([Environment]::CommandLine -match '-NonI') { return }

$dotfilesDir = if ($env:DOTFILES_DIR) { $env:DOTFILES_DIR } else { "$HOME/.dotfiles" }
if (Test-Path "$dotfilesDir/.path.win") {
    Get-Content "$dotfilesDir/.path.win" | Where-Object { $_ } | ForEach-Object {
        if (($env:PATH -split ';') -notcontains $_) {
            $env:PATH += ";$_"
        }
    }
}

pwshac | Out-String -Width ([int]::MaxValue) | Invoke-Expression

dotsh pwsh (mise env --dotenv) | Out-String -Width ([int]::MaxValue) | Invoke-Expression

if (Test-Path "$dotfilesDir/.env") { dotsh pwsh (Get-Content "$dotfilesDir/.env" -Raw) | Out-String -Width ([int]::MaxValue) | Invoke-Expression }
if (Test-Path "$dotfilesDir/.local/.env") { dotsh pwsh (Get-Content "$dotfilesDir/.local/.env" -Raw) | Out-String -Width ([int]::MaxValue) | Invoke-Expression }

direnv hook pwsh | Out-String -Width ([int]::MaxValue) | Invoke-Expression

shaka pwsh | Out-String -Width ([int]::MaxValue) | Invoke-Expression
zoxide init powershell | Out-String -Width ([int]::MaxValue) | Invoke-Expression
zoxide add "$PWD" 2>$null
starship init powershell | Out-String -Width ([int]::MaxValue) | Invoke-Expression
