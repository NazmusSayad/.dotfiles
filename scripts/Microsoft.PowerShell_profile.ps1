# Alias
Set-Alias -Name y -Value yarn
Set-Alias -Name live -Value live-server
function ya {iex "yarn add $($args)"}  
function yad {iex "yarn add -D $($args)"}
function n {iex "node $($args)"}
function nw {iex "node --watch $($args)"}
function tn {iex "ts-node $($args)"}
function tnw {iex "ts-node-dev $($args)"}
function nm {iex "nodemon $($args)"}

# Functions
function ghp {
  $msg=($args -join " ").Trim()
  if ([string]::IsNullOrEmpty($msg)) {
    Write-Host "You must need to give a <commit_message>" -ForegroundColor Red
    return
  }
  
  git status --short
  Write-Host ""

  git add -A | Out-Null
  git commit -m "$msg"
  Write-Host ""

  git push --quiet
}

function ghr {
  param (
    [string]$hash
  )

  if ([string]::IsNullOrEmpty($hash)) {
    Write-Host "You must need to give a <commit_id>" -ForegroundColor Red
    return
  }

  git checkout $hash . 
  ghp ("$hash restored; " + ($args[1..($args.Length - 1)] -join " "))
} 

function ghhr {
  param (
    [string]$branch = "master"
  )

  git checkout --orphan latest_branch
  git add -A
  git commit -am "initial commit"
  git branch -D $branch
  git branch -m $branch
  git push -f origin $branch
}

function ghc {
  param (
    [string]$repo
  )

  if ($repo -match "^http.*") {
    git clone $repo
  }
  else {
    gh repo clone $repo
  }
}
