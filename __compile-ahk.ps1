echo "> Deleting old compiled files..."
Remove-Item .\build\ahk -Recurse -Force -ErrorAction SilentlyContinue

echo ""
echo "> Compiling ahk scripts..."
& ".\build\bin\ahk-compile.exe"

echo ""
Read-Host "Press Enter to continue..."

