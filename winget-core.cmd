@echo off
setlocal enabledelayedexpansion

set PACKAGES=

for /f "usebackq delims=" %%a in ("winget-apps.ini") do (
  set line=%%a
  if not "!line!"=="" if not "!line:~0,1!"=="#" (
    for /f "tokens=1" %%b in ("!line!") do (
      set package=%%b
      if not "!package!"=="" set PACKAGES=!PACKAGES! !package!
    )
  )
)

endlocal
echo !PACKAGES!
for /f "tokens=* delims= " %%a in ("!PACKAGES!") do set PACKAGES=%%a
