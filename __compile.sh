#!/bin/bash
set -e

echo "> Killing all AHK scripts..."
killed_ahk=()
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || "$OSTYPE" == "win32" ]]; then
  while IFS= read -r name; do
    name="${name%.exe}"
    [[ -n "$name" ]] && killed_ahk+=("$name")
  done < <(tasklist //FI "IMAGENAME eq AHK-*" 2>/dev/null | awk 'NR>=3 && /AHK-/{print $1}' | sort -u)
  if [ ${#killed_ahk[@]} -gt 0 ]; then
    sudo taskkill //F //IM AHK-* 2>/dev/null || true
  fi
fi

echo "> Cleaning build directory..."
rm -rf ./.build/bin ./.build/ahk

echo ""
echo "> Compiling AutoHotkey scripts..."
go run ./src/compile-ahk/main.go

echo ""
echo "> Compiling Go scripts..."
go run ./src/compile-scripts/main.go

if [ ${#killed_ahk[@]} -gt 0 ]; then
  echo ""
  echo "> Restarting AHK scripts..."
  ahk_dir="./.build/ahk"
  for name in "${killed_ahk[@]}"; do
    exe="$ahk_dir/$name.exe"
    echo "> Restarting: $exe"
    if [ -f "$exe" ]; then
      powershell.exe -NoProfile -Command "Start-Process '$exe' -Verb RunAs"
    fi
  done
fi

echo ""
echo "Done!"
