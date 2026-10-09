$ErrorActionPreference = 'Stop'

$keepAlive = Get-CimInstance Win32_Process | Where-Object {
    $_.Name -eq 'wsl.exe' -and $_.CommandLine -like '*Ubuntu-26.04*' -and $_.CommandLine -like '*/bin/sleep infinity*'
} | Select-Object -First 1

if ($null -eq $keepAlive) {
    Start-Process wsl.exe -ArgumentList '-d Ubuntu-26.04 --cd ~ --exec /bin/sleep infinity' -WindowStyle Hidden
}

wsl.exe -d Ubuntu-26.04 -u root --exec /bin/systemctl start ttysh
if ($LASTEXITCODE -ne 0) { throw 'Failed to start ttysh in Ubuntu.' }
