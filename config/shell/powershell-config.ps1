$env:GOBIN = (go env GOBIN)
$env:GOROOT = (go env GOROOT)
$env:JAVA_HOME = (mise where java)

Remove-Item Alias:ni -Force -ErrorAction Ignore

function ni {
  & flutter pub get $args
}
function nid {
  & flutter pub get -d $args
}

function fl {
  & flutter $args
}
function flp {
  & flutter pub $args
}
function flr {
  & flutter run $args
}
function fle {
  & flutter emulators $args
}

function d {
  & docker $args
}
function dc {
  & docker compose $args
}
function dcu {
  & docker compose up -d $args
}
function dcd {
  & docker compose down $args
}

Invoke-Expression "$(direnv hook pwsh)"
Invoke-Expression (&starship init powershell)
Invoke-Expression (& { (zoxide init powershell | Out-String) })
