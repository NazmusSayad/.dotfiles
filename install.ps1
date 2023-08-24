$PROFILEDIR="D:\$env:USERNAME\Documents\WindowsPowerShell"
$isDirExists=Test-Path $PROFILEDIR
if($isDirExists -eq $False){
  md $PROFILEDIR
}
Copy-Item "./scripts/Microsoft.PowerShell_profile.ps1" $PROFILE -Force
echo "Powershell profile copied..."

git config --global core.excludesfile %USERPROFILE%\.gitignore

Rename-Item -Path "./config" -NewName "./config.ps1"
powershell.exe -File config.ps1
Rename-Item -Path "./config.ps1" -NewName "./config"

Write-Host 'Press any key to exit...' -NoNewLine -ForegroundColor Red
$null = $Host.UI.RawUI.ReadKey('NoEcho,IncludeKeyDown')
