Write-Output 'Setting Windows Settings...'

# Removes OneDrive
Remove-Item "C:\Users\Default\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\OneDrive.lnk" -ErrorAction Continue
Remove-Item "C:\Users\Default\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\OneDrive.exe" -ErrorAction Continue
Remove-Item "C:\Windows\System32\OneDriveSetup.exe" -ErrorAction Continue
Remove-Item "C:\Windows\SysWOW64\OneDriveSetup.exe" -ErrorAction Continue
