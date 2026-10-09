$ErrorActionPreference = 'Continue'

while ($true) {
    wsl.exe -d Ubuntu-26.04 --cd /home/sayad --exec /bin/sleep infinity
    Start-Sleep -Seconds 2
}
