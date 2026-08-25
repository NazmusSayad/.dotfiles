if ([Environment]::CommandLine -match '-NonI') { return }

if (Test-Path ~/.path) {
    Get-Content ~/.path | Where-Object { $_ } | ForEach-Object {
        if (($env:PATH -split ';') -notcontains $_) {
            $env:PATH += ";$_"
        }
    }
}

pwshac | Out-String -Width ([int]::MaxValue) | Invoke-Expression

dotsh pwsh (mise env --dotenv) | Out-String -Width ([int]::MaxValue) | Invoke-Expression

if (Test-Path ~/.dotfiles/.env) { dotsh pwsh (Get-Content ~/.dotfiles/.env -Raw) | Out-String -Width ([int]::MaxValue) | Invoke-Expression }
if (Test-Path ~/.env) { dotsh pwsh (Get-Content ~/.env -Raw) | Out-String -Width ([int]::MaxValue) | Invoke-Expression }

direnv hook pwsh | Out-String -Width ([int]::MaxValue) | Invoke-Expression

shaka pwsh | Out-String -Width ([int]::MaxValue) | Invoke-Expression
zoxide init powershell | Out-String -Width ([int]::MaxValue) | Invoke-Expression
starship init powershell | Out-String -Width ([int]::MaxValue) | Invoke-Expression
