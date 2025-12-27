Write-Output 'Removing Windows Features...'

$featuresToDisable = @(
  'Microsoft-SnippingTool',
  'Microsoft-Windows-Hello-Face',
  'Print-Fax-Client'
)

Get-WindowsOptionalFeature -Online |
ForEach-Object {
  $featureName = $_.FeatureName
  if ($featuresToDisable | Where-Object { $featureName -Like $_ }) {
    Write-Output "Removing $featureName..."
    Disable-WindowsOptionalFeature -Online -FeatureName $featureName -NoRestart
  }
}