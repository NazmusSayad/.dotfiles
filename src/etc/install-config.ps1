Write-Host "> Setting up git config..."
bash "./src/shell/git-config.sh"

Write-Host ""
Write-Host "> Symlinking..."
go run ./src/scripts/symlink-init/main.go

Write-Host ""
Write-Host "> Installing tasks..."
go run ./src/install-windows-tasks/main.go

Write-Host ""
Write-Host "> Installing start menu entries..."
go run ./src/install-start-menu/main.go

Write-Host ""
Write-Host "Done!"