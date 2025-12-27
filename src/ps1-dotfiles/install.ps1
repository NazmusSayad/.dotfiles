$CURRENT_DIR = Split-Path -LiteralPath $MyInvocation.MyCommand.Path -Parent
$CURRENT_DIR = $CURRENT_DIR.TrimEnd('\')

Write-Host "Setting DOTFILES_DIR to $CURRENT_DIR ..."
# [Environment]::SetEnvironmentVariable('DOTFILES_DIR', $CURRENT_DIR, 'User')

$userPath = [Environment]::GetEnvironmentVariable('PATH', 'User')
Write-Host "Setting PATH to $CURRENT_DIR\.build\bin;$userPath ..."
# [Environment]::SetEnvironmentVariable('PATH', "$CURRENT_DIR\.build\bin;$userPath", 'User')

Write-Host
Write-Host 'SUCCESS: Dotfiles installation completed.'
cmd /c pause
