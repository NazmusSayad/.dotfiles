$ErrorActionPreference = 'Stop'

$keepAliveScript = Join-Path $PSScriptRoot 'keep-alive.ps1'
$keepAlive = Get-CimInstance Win32_Process | Where-Object {
    $_.Name -eq 'powershell.exe' -and $_.CommandLine -like "*$keepAliveScript*"
} | Select-Object -First 1

if ($null -eq $keepAlive) {
    Start-Process powershell.exe -ArgumentList @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', $keepAliveScript) -WindowStyle Hidden
}

wsl.exe -d Ubuntu-26.04 -u root --exec /bin/systemctl start sshtty
if ($LASTEXITCODE -ne 0) { throw 'Failed to start sshtty in Ubuntu.' }

$server = Get-CimInstance Win32_Process | Where-Object {
    $_.Name -eq 'sshtty.exe' -and $_.CommandLine -match '--port\s+47831\b'
} | Select-Object -First 1

if ($null -eq $server) {
    $mise = (Get-Command mise -ErrorAction Stop).Source
    Start-Process $mise -ArgumentList 'exec -- sshtty --host 127.0.0.1 --port 47831' -WorkingDirectory $HOME -WindowStyle Hidden
}

$tunnel = Get-CimInstance Win32_Process | Where-Object {
    $_.Name -eq 'cloudflared.exe' -and $_.CommandLine -match 'tunnel\s+run'
} | Select-Object -First 1

if ($null -eq $tunnel) {
    $logDirectory = Join-Path $env:LOCALAPPDATA 'Temp\desktop-tunnel'
    New-Item -ItemType Directory -Path $logDirectory -Force | Out-Null
    $cloudflared = (Get-Command cloudflared -ErrorAction Stop).Source
    Start-Process $cloudflared -ArgumentList 'tunnel run' -WindowStyle Hidden -RedirectStandardOutput (Join-Path $logDirectory 'stdout.log') -RedirectStandardError (Join-Path $logDirectory 'stderr.log')
}
