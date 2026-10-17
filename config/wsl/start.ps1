$ErrorActionPreference = 'Stop'

$keepAliveScript = Join-Path $PSScriptRoot 'keep-alive.ps1'
$keepAlive = Get-CimInstance Win32_Process | Where-Object {
    $_.Name -eq 'powershell.exe' -and $_.CommandLine -like "*$keepAliveScript*"
} | Select-Object -First 1

if ($null -eq $keepAlive) {
    Start-Process powershell.exe -ArgumentList @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', $keepAliveScript) -WindowStyle Hidden
}

wsl.exe -d Ubuntu-26.04 -u root --exec /bin/systemctl start offsite
if ($LASTEXITCODE -ne 0) { throw 'Failed to start offsite in Ubuntu.' }

$mise = (Get-Command mise -ErrorAction Stop).Source

$server = Get-CimInstance Win32_Process | Where-Object {
    $_.Name -eq 'sshtty.exe' -and $_.CommandLine -match '--port\s+20203\b'
} | Select-Object -First 1

if ($null -eq $server) {
    Start-Process $mise -ArgumentList 'exec -- sshtty --host 127.0.0.1 --port 20203' -WorkingDirectory $HOME -WindowStyle Hidden
}
