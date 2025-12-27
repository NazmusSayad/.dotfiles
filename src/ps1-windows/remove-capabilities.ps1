Write-Output 'Removing Windows Capabilities...'

$capabilitiesToRemove = @(
  'Browser.InternetExplorer',
  'MathRecognizer',
  'OpenSSH.Client',
  'Microsoft.Windows.MSPaint',
  'Microsoft.Windows.PowerShell.ISE',
  'App.Support.QuickAssist',
  'App.StepsRecorder',
  'Media.WindowsMediaPlayer',
  'Microsoft.Windows.WordPad'
)

Get-WindowsCapability -Online |
ForEach-Object {
  $capabilityName = $_.Name
  if ($capabilitiesToRemove | Where-Object { $capabilityName -Like $_ }) {
    Write-Output "Removing $capabilityName..."
    Remove-WindowsCapability -Online -Name $capabilityName
  }
}