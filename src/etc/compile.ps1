Write-Host "> Killing all AHK scripts..."
$ahk = Get-Process -ErrorAction SilentlyContinue | Where-Object { $_.Name -like "AHK-*" }
$killedAhk = @()
if ($ahk) {
    $killedAhk = @($ahk.Name | Sort-Object -Unique)
    sudo taskkill /F /IM AHK-*
}

Write-Host "> Cleaning build directory..."
Remove-Item .\.build\bin -Recurse -Force -ErrorAction SilentlyContinue
Remove-Item .\.build\ahk -Recurse -Force -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "> Compiling AutoHotkey scripts..."
go run ./src/compile-ahk/main.go

Write-Host ""
Write-Host "> Compiling aliases..."
go run ./src/compile-alias/main.go

Write-Host ""
Write-Host "> Compiling Go scripts..."
go run ./src/compile-scripts/main.go

if ($killedAhk.Count -gt 0) {
    Write-Host ""
    Write-Host "> Restarting AHK scripts..."
    $ahkDir = Join-Path $PWD ".build\ahk"
    foreach ($name in $killedAhk) {
        $exe = Join-Path $ahkDir "$name.exe"
        Write-Host "> Restarting: $exe"
        if (Test-Path $exe) {
            Start-Process $exe -Verb RunAs
        }
    }
}
