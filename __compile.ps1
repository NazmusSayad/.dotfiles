echo "> Deleting old compiled files..."
Remove-Item .\build\bin -Recurse -Force -ErrorAction SilentlyContinue

$goTargets = @{
    ".\src\startup\main.go" = ".\build\bin\windows-startup.exe"
    ".\src\vbproxy\main.go" = ".\build\bin\vbproxy.exe"
    ".\src\ahk\main.go" = ".\build\bin\ahk-compile.exe"
    ".\src\msys-setup\main.go" = ".\build\bin\msys-setup.exe"
    ".\src\symlink\main.go" = ".\build\bin\symlink-config.exe"
    ".\src\clean-code-snippets\main.go" = ".\build\bin\clean-code-snippets.exe"
    ".\src\winget\install\main.go" = ".\build\bin\winget-install.exe"
    ".\src\winget\upgrade\main.go" = ".\build\bin\winget-upgrade.exe"
    ".\src\slack\status\main.go" = ".\build\bin\slack-status.exe"
    ".\src\slack\startup\main.go" = ".\build\bin\slack-startup.exe"
}

foreach ($entry in $goTargets.GetEnumerator()) {
    echo ""
    echo "> Compiling $($entry.Key) (Go)..."
    go build -o $entry.Value $entry.Key
}

$rustTargets = @{
    ".\src\functions\c.rs" = ".\build\bin\c.exe"
    ".\src\functions\gc.rs" = ".\build\bin\gc.exe"
    ".\src\functions\gr.rs" = ".\build\bin\gr.exe"
    ".\src\functions\gback.rs" = ".\build\bin\gback.exe"
    ".\src\functions\gp.rs" = ".\build\bin\gp.exe"
    ".\src\functions\gpr.rs" = ".\build\bin\gpr.exe"
    ".\src\functions\gclean.rs" = ".\build\bin\gclean.exe"
    ".\src\functions\greset.rs" = ".\build\bin\greset.exe"
    ".\src\functions\gpg_unlock.rs" = ".\build\bin\gpg-unlock.exe"
    ".\src\functions\fscs.rs" = ".\build\bin\fscs.exe"
}

foreach ($entry in $rustTargets.GetEnumerator()) {
    echo ""
    echo "> Compiling $($entry.Key) (Rust)..."
    rustc -C strip=symbols -Clink-arg=/DEBUG:NONE -Clink-arg=/PDB:NONE $entry.Key -o $entry.Value
}

echo ""
Read-Host "Press Enter to continue..."
