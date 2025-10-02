Write-Output 'Removing Windows Apps...'

$appxPackagesToRemove = @(
  'Microsoft.Wallet',
  'Microsoft.Windows.DevHome',
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
  'Microsoft.BingSearch',
  'Microsoft.WindowsCamera',
  'Microsoft.Windows.Photos',
  'Microsoft.WindowsCalculator',
  'Microsoft.Windows.DevHome',
  'Clipchamp.Clipchamp',
  'Microsoft.WindowsAlarms',
  'Microsoft.549981C3F5F10',
  'MicrosoftCorporationII.MicrosoftFamily',
  'Microsoft.WindowsFeedbackHub',
  'Microsoft.GetHelp', 
  'microsoft.windowscommunicationsapps',
  'Microsoft.WindowsMaps',
  'Microsoft.ZuneVideo',
  'Microsoft.BingNews',
  'Microsoft.MicrosoftOfficeHub',
  'Microsoft.Office.OneNote',
  'Microsoft.OutlookForWindows',
  'Microsoft.Paint',
  'Microsoft.MSPaint',
  'Microsoft.People',
  'Microsoft.PowerAutomateDesktop',
  'MicrosoftCorporationII.QuickAssist',
  'Microsoft.SkypeApp',
  'Microsoft.MicrosoftSolitaireCollection',
  'Microsoft.MicrosoftStickyNotes',
  'MSTeams',
  'Microsoft.Getstarted',
  'Microsoft.Todos',
  'Microsoft.WindowsSoundRecorder',
  'Microsoft.BingWeather',
  'Microsoft.ZuneMusic',
  'Microsoft.Xbox*',
  'Microsoft.GamingApp',
  'Microsoft.YourPhone',
  'Microsoft.MicrosoftEdge*',
  'Microsoft.OneDrive',
  'Microsoft.549981C3F5F10',
  'Microsoft.MixedReality.Portal',
  'Microsoft.Windows.Ai.Copilot.Provider',
  'Microsoft.WindowsMeetNow',
  'Microsoft.WindowsStore',
  'Microsoft.ScreenSketch',
  'Microsoft.SnippingTool',
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
