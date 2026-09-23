Write-Output 'Removing Windows Apps...'

$appxPackagesToRemove = @(
  'Microsoft.Wallet',
  'Microsoft.StorePurchaseApp',
  'Microsoft.BioEnrollment',
  'Microsoft.Windows.CloudExperienceHost',
  'Microsoft.Windows.ContentDeliveryManager',
  'Microsoft.Windows.PeopleExperienceHost',
  'Microsoft.Windows.OOBENetworkCaptivePortal',
  'Microsoft.Windows.OOBENetworkConnectionFlow',
  'Microsoft.Windows.CapturePicker',
  'Microsoft.Windows.SecureAssessmentBrowser',
  'Microsoft.MicrosoftEdgeDevToolsClient',
  'Microsoft.Windows.XGpuEjectDialog',
  'Microsoft.XboxGameCallableUI',
  'NcsiUwpApp'
)

$packagesToRemove = @(
  'Microsoft.Microsoft3DViewer',
  'Microsoft.Print3D',
  'Microsoft.MSPaint',
  'Microsoft.WindowsCamera',
  'Microsoft.Windows.Photos',
  'Microsoft.WindowsCalculator',
  'Microsoft.WindowsAlarms',
  'Microsoft.549981C3F5F10',
  'Microsoft.WindowsFeedbackHub',
  'Microsoft.GetHelp',
  'microsoft.windowscommunicationsapps',
  'Microsoft.WindowsMaps',
  'Microsoft.ZuneVideo',
  'Microsoft.BingNews',
  'Microsoft.BingWeather',
  'Microsoft.MicrosoftOfficeHub',
  'Microsoft.Office.OneNote',
  'Microsoft.People',
  'Microsoft.Messaging',
  'Microsoft.OneConnect',
  'Microsoft.SkypeApp',
  'Microsoft.MicrosoftSolitaireCollection',
  'Microsoft.MicrosoftStickyNotes',
  'Microsoft.Getstarted',
  'Microsoft.Todos',
  'Microsoft.WindowsSoundRecorder',
  'Microsoft.ZuneMusic',
  'Microsoft.Xbox*',
  'Microsoft.GamingApp',
  'Microsoft.YourPhone',
  'Microsoft.MicrosoftEdge*',
  'Microsoft.OneDrive',
  'Microsoft.MixedReality.Portal',
  'Microsoft.WindowsStore',
  'Microsoft.ScreenSketch',
  'Microsoft.XboxGameCallableUI',
  'Microsoft.Windows.NarratorQuickStart',
  'Microsoft.Windows.PeopleExperienceHost',
  'Microsoft.Windows.ParentalControls',
  'Microsoft.Windows.CloudExperienceHost',
  'Microsoft.MicrosoftEdgeDevToolsClient',
  'AppUp.IntelGraphicsExperience'
)

$packagesToRemove += $appxPackagesToRemove

Get-AppxProvisionedPackage -Online |
ForEach-Object {
  $packageName = $_.DisplayName
  if ($packagesToRemove | Where-Object { $packageName -Like $_ }) {
    Write-Output "Removing $packageName..."
    Remove-AppxProvisionedPackage -AllUsers -Online -PackageName $_.PackageName
  }
}

Get-AppxPackage |
ForEach-Object {
  $packageName = $_.Name
  if ($appxPackagesToRemove | Where-Object { $packageName -Like $_ }) {
    Write-Output "Removing $packageName..."
    Remove-AppxPackage -Package $_.PackageFullName
  }
}
