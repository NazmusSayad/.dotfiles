Write-Output 'Applying Gaming Performance Tweaks...'

# POWER PLAN
# Enables and Activates the Ultimate Performance Power Plan
$guidPattern = '[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}'
$ultimate = powercfg /list | Select-String 'Ultimate Performance'
if (-not $ultimate) {
    $ultimate = powercfg -duplicatescheme e9a42b02-d5df-448d-aa00-03f14749eb61
}
if (($ultimate | Out-String) -match $guidPattern) {
    powercfg /setactive $Matches[0]
} else {
    # Falls back to High Performance
    powercfg /setactive 8c5e7fda-e8bf-4a96-9a85-a6e23a8c635c
}

# Keeps the CPU at Full Speed and Disables Core Parking
powercfg /setacvalueindex scheme_current sub_processor PROCTHROTTLEMIN 100
powercfg /setacvalueindex scheme_current sub_processor PROCTHROTTLEMAX 100
powercfg /setacvalueindex scheme_current sub_processor CPMINCORES 100
# Disables USB Selective Suspend (Prevents Mouse / Keyboard / Controller Wake Lag)
powercfg /setacvalueindex scheme_current 2a737441-1930-4402-8d77-b2bebba308a3 48e6b7a6-50f5-4782-a5d4-53bb8f07e226 0
# Disables PCI Express Link State Power Management (GPU)
powercfg /setacvalueindex scheme_current sub_pciexpress ASPM 0
# Never Turn Off Disks or Sleep
powercfg /setacvalueindex scheme_current sub_disk DISKIDLE 0
powercfg /change standby-timeout-ac 0
powercfg /setactive scheme_current

# Disables Hibernation (Also Disables Fast Startup, Frees Disk Space)
powercfg /hibernate off

# Disables Power Throttling for Background Processes
reg.exe add "HKLM\SYSTEM\CurrentControlSet\Control\Power\PowerThrottling" /v PowerThrottlingOff /t REG_DWORD /d 1 /f

# CPU / GPU SCHEDULING
# Prioritizes the Foreground Application (Short, Fixed Quantum, High Foreground Boost)
reg.exe add "HKLM\SYSTEM\CurrentControlSet\Control\PriorityControl" /v Win32PrioritySeparation /t REG_DWORD /d 38 /f

# Enables Hardware Accelerated GPU Scheduling (Windows 10 2004+ and Supported GPU)
reg.exe add "HKLM\SYSTEM\CurrentControlSet\Control\GraphicsDrivers" /v HwSchMode /t REG_DWORD /d 2 /f

# Disables Virtualization Based Security and Memory Integrity (HVCI)
reg.exe add "HKLM\SYSTEM\CurrentControlSet\Control\DeviceGuard" /v EnableVirtualizationBasedSecurity /t REG_DWORD /d 0 /f
reg.exe add "HKLM\SYSTEM\CurrentControlSet\Control\DeviceGuard\Scenarios\HypervisorEnforcedCodeIntegrity" /v Enabled /t REG_DWORD /d 0 /f

# GAME MODE / GAME BAR
# Enables Game Mode
reg.exe add "HKEY_CURRENT_USER\SOFTWARE\Microsoft\GameBar" /v AutoGameModeEnabled /t REG_DWORD /d 1 /f
reg.exe add "HKEY_CURRENT_USER\SOFTWARE\Microsoft\GameBar" /v AllowAutoGameMode /t REG_DWORD /d 1 /f

# Disables Game Bar and Background Recording
reg.exe add "HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows\CurrentVersion\GameDVR" /v AppCaptureEnabled /t REG_DWORD /d 0 /f
reg.exe add "HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows\CurrentVersion\GameDVR" /v HistoricalCaptureEnabled /t REG_DWORD /d 0 /f
reg.exe add "HKEY_CURRENT_USER\SOFTWARE\Microsoft\GameBar" /v ShowStartupPanel /t REG_DWORD /d 0 /f

# INPUT
# Disables Mouse Acceleration (Enhance Pointer Precision)
reg.exe add "HKEY_CURRENT_USER\Control Panel\Mouse" /v MouseSpeed /t REG_SZ /d 0 /f
reg.exe add "HKEY_CURRENT_USER\Control Panel\Mouse" /v MouseThreshold1 /t REG_SZ /d 0 /f
reg.exe add "HKEY_CURRENT_USER\Control Panel\Mouse" /v MouseThreshold2 /t REG_SZ /d 0 /f

# Fastest Keyboard Repeat Rate and Shortest Delay
reg.exe add "HKEY_CURRENT_USER\Control Panel\Keyboard" /v KeyboardDelay /t REG_SZ /d 0 /f
reg.exe add "HKEY_CURRENT_USER\Control Panel\Keyboard" /v KeyboardSpeed /t REG_SZ /d 31 /f

# VISUAL EFFECTS (System Properties > Advanced > Performance Options)
# Custom: Everything Off Except "Smooth edges of screen fonts" and "Show thumbnails instead of icons"
reg.exe add "HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\VisualEffects" /v VisualFXSetting /t REG_DWORD /d 3 /f
# Animate controls, Fade/Slide menus, Fade/Slide tooltips, Fade out menu items, Mouse pointer shadow,
# Window shadows, Slide open combo boxes, Smooth-scroll list boxes
reg.exe add "HKEY_CURRENT_USER\Control Panel\Desktop" /v UserPreferencesMask /t REG_BINARY /d 9012038010000000 /f
# Animate windows when minimizing and maximizing
reg.exe add "HKEY_CURRENT_USER\Control Panel\Desktop\WindowMetrics" /v MinAnimate /t REG_SZ /d 0 /f
# Animations in the taskbar
reg.exe add "HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Advanced" /v TaskbarAnimations /t REG_DWORD /d 0 /f
# Enable Peek
reg.exe add "HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows\DWM" /v EnableAeroPeek /t REG_DWORD /d 0 /f
# Save taskbar thumbnail previews
reg.exe add "HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows\DWM" /v AlwaysHibernateThumbnails /t REG_DWORD /d 0 /f
# Show thumbnails instead of icons (Kept On)
reg.exe add "HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Advanced" /v IconsOnly /t REG_DWORD /d 0 /f
# Show translucent selection rectangle
reg.exe add "HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Advanced" /v ListviewAlphaSelect /t REG_DWORD /d 0 /f
# Show window contents while dragging
reg.exe add "HKEY_CURRENT_USER\Control Panel\Desktop" /v DragFullWindows /t REG_SZ /d 0 /f
# Smooth edges of screen fonts (Kept On)
reg.exe add "HKEY_CURRENT_USER\Control Panel\Desktop" /v FontSmoothing /t REG_SZ /d 2 /f
# Use drop shadows for icon labels on the desktop
reg.exe add "HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Advanced" /v ListviewShadow /t REG_DWORD /d 0 /f
# Removes the Delay Before Menus Open
reg.exe add "HKEY_CURRENT_USER\Control Panel\Desktop" /v MenuShowDelay /t REG_SZ /d 0 /f

# BACKGROUND ACTIVITY
# Disables Background Apps
reg.exe add "HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows\CurrentVersion\BackgroundAccessApplications" /v GlobalUserDisabled /t REG_DWORD /d 1 /f
reg.exe add "HKLM\SOFTWARE\Policies\Microsoft\Windows\AppPrivacy" /v LetAppsRunInBackground /t REG_DWORD /d 2 /f

# Disables Automatic Maintenance (Defrag, Diagnostics Running Mid-Game)
reg.exe add "HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Schedule\Maintenance" /v MaintenanceDisabled /t REG_DWORD /d 1 /f

# NETWORK
# Disables Nagle's Algorithm on All Network Interfaces (Lower Online Latency)
Get-ChildItem "HKLM:\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces" |
ForEach-Object {
    $interfacePath = $_.Name
    reg.exe add "$interfacePath" /v TcpAckFrequency /t REG_DWORD /d 1 /f
    reg.exe add "$interfacePath" /v TCPNoDelay /t REG_DWORD /d 1 /f
}
