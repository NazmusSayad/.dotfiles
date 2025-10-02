Write-Output 'Setting Windows Settings...'

# Removes OneDrive
Remove-Item "C:\Users\Default\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\OneDrive.lnk" -ErrorAction Continue
Remove-Item "C:\Users\Default\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\OneDrive.exe" -ErrorAction Continue
Remove-Item "C:\Windows\System32\OneDriveSetup.exe" -ErrorAction Continue
Remove-Item "C:\Windows\SysWOW64\OneDriveSetup.exe" -ErrorAction Continue

# Enable the Ultimate Performance power plan
powercfg.exe -duplicatescheme e9a42b02-d5df-448d-aa00-03f14749eb61
powercfg.exe -setactive e9a42b02-d5df-448d-aa00-03f14749eb61