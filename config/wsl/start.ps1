$ErrorActionPreference = 'Stop'

$keepAlive = Get-CimInstance Win32_Process | Where-Object {
    $_.Name -eq 'wsl.exe' -and $_.CommandLine -like '*Ubuntu-26.04*' -and $_.CommandLine -like '*/bin/sleep infinity*'
} | Select-Object -First 1

if ($null -eq $keepAlive) {
    Start-Process wsl.exe -ArgumentList '-d Ubuntu-26.04 --cd ~ --exec /bin/sleep infinity' -WindowStyle Hidden
}

wsl.exe -d Ubuntu-26.04 -u root --exec /bin/systemctl start ttysh
if ($LASTEXITCODE -ne 0) { throw 'Failed to start ttysh in Ubuntu.' }

$tunnel = Get-CimInstance Win32_Process | Where-Object {
    $_.Name -eq 'cloudflared.exe' -and $_.CommandLine -match 'tunnel\s+run'
} | Select-Object -First 1

if ($null -eq $tunnel) {
    $logDirectory = Join-Path $env:LOCALAPPDATA 'Temp\desktop-tunnel'
    New-Item -ItemType Directory -Path $logDirectory -Force | Out-Null
    $cloudflared = (Get-Command cloudflared -ErrorAction Stop).Source
    Start-Process $cloudflared -ArgumentList 'tunnel run' -WindowStyle Hidden -RedirectStandardOutput (Join-Path $logDirectory 'stdout.log') -RedirectStandardError (Join-Path $logDirectory 'stderr.log')
}
