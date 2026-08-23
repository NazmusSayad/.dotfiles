if ([Environment]::CommandLine -match '-NonI') { return }

if (Test-Path ~/.path) {
    Get-Content ~/.path | Where-Object { $_ } | ForEach-Object {
        if (($env:PATH -split ';') -notcontains $_) {
            $env:PATH += ";$_"
        }
    }
}

Invoke-Expression (pwshac cd ls cp mv rm cat echo pwd man kill sleep ps history clear | Out-String -Width ([int]::MaxValue))

dotsh pwsh (mise env --dotenv) | Out-String -Width ([int]::MaxValue) | Where-Object { $_.Trim() } | Invoke-Expression

if (Test-Path ~/.dotfiles/.env) { dotsh pwsh (Get-Content ~/.dotfiles/.env -Raw) | Out-String -Width ([int]::MaxValue) | Where-Object { $_.Trim() } | Invoke-Expression }
if (Test-Path ~/.env) { dotsh pwsh (Get-Content ~/.env -Raw) | Out-String -Width ([int]::MaxValue) | Where-Object { $_.Trim() } | Invoke-Expression }

direnv hook pwsh | Out-String -Width ([int]::MaxValue) | Where-Object { $_.Trim() } | Invoke-Expression

shaka pwsh | Out-String -Width ([int]::MaxValue) | Where-Object { $_.Trim() } | Invoke-Expression
zoxide init powershell | Out-String -Width ([int]::MaxValue) | Where-Object { $_.Trim() } | Invoke-Expression
starship init powershell | Out-String -Width ([int]::MaxValue) | Where-Object { $_.Trim() } | Invoke-Expression
