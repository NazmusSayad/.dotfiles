# Disables Telemetry Scheduled Tasks
$scheduledTasks = @(
    "Microsoft\Windows\Application Experience\Microsoft Compatibility Appraiser",
    "Microsoft\Windows\Application Experience\ProgramDataUpdater",
    "Microsoft\Windows\Application Experience\StartupAppTask",
    "Microsoft\Windows\Autochk\Proxy",
    "Microsoft\Windows\Customer Experience Improvement Program\Consolidator",
    "Microsoft\Windows\Customer Experience Improvement Program\KernelCeipTask",
    "Microsoft\Windows\Customer Experience Improvement Program\UsbCeip",
    "Microsoft\Windows\DiskDiagnostic\Microsoft-Windows-DiskDiagnosticDataCollector",
    "Microsoft\Windows\Feedback\Siuf\DmClient",
    "Microsoft\Windows\Feedback\Siuf\DmClientOnScenarioDownload",
    "Microsoft\Windows\Windows Error Reporting\QueueReporting",
    "Microsoft\Windows\CloudExperienceHost\CreateObjectTask",
    "Microsoft\Windows\Maps\MapsToastTask",
    "Microsoft\Windows\Maps\MapsUpdateTask"
)
foreach ($task in $scheduledTasks) {
    schtasks /Change /TN $task /Disable
}
