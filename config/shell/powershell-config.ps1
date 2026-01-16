Remove-Item Alias:ni -Force -ErrorAction Ignore
Invoke-Expression (&starship init powershell)
Invoke-Expression (& { (zoxide init powershell | Out-String) })

$env:GOPATH = (go env GOPATH)
$env:GOROOT = (go env GOROOT)
$env:JAVA_HOME = (mise where java)