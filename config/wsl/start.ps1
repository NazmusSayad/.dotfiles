$ErrorActionPreference = 'Stop'

$keepAliveScript = Join-Path $PSScriptRoot 'keep-alive.ps1'
$keepAlive = Get-CimInstance Win32_Process | Where-Object {
    $_.Name -eq 'powershell.exe' -and $_.CommandLine -like "*$keepAliveScript*"
} | Select-Object -First 1

if ($null -eq $keepAlive) {
    Start-Process powershell.exe -ArgumentList @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', $keepAliveScript) -WindowStyle Hidden
}

wsl.exe -d Ubuntu-26.04 -u root --exec /bin/systemctl start offsite sshtty
if ($LASTEXITCODE -ne 0) { throw 'Failed to start offsite and sshtty in Ubuntu.' }

$mise = (Get-Command mise -ErrorAction Stop).Source

$offsite = Get-CimInstance Win32_Process | Where-Object {
    $_.Name -eq 'node.exe' -and $_.CommandLine -like '*offsite.js*' -and $_.CommandLine -match '--port\s+20202\b'
} | Select-Object -First 1

if ($null -eq $offsite) {
    Start-Process $mise -ArgumentList 'exec -- offsite serve --no-auth --no-pair --host 127.0.0.1 --port 20202' -WorkingDirectory $HOME -WindowStyle Hidden
}

$server = Get-CimInstance Win32_Process | Where-Object {
    $_.Name -eq 'sshtty.exe' -and $_.CommandLine -match '--port\s+20203\b'
} | Select-Object -First 1

if ($null -eq $server) {
    Start-Process $mise -ArgumentList 'exec -- sshtty --host 127.0.0.1 --port 20203' -WorkingDirectory $HOME -WindowStyle Hidden
}
