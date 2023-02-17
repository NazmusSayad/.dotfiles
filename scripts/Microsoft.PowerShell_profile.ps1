Function ghc{param($p);gh repo clone $p}
Function ghp{param($p);if($p.Length -eq 0){$p=Read-Host "#"}if($p -eq ''){$p="..."}git status;git add .;git commit -m $p;git push}
Function ghi{Remove-Item .git -Force -Recurse;cls;git init;git add .;git commit -m "init..."}
Function df{param($p);If(([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]"Administrator")){Defrag.exe /C;if($p -eq 'r'){shutdown.exe -r -t 0}elseif($p -eq 's'){shutdown.exe -s -t 0}}}
Function ghhr{param($branch);if($branch.Length -eq 0){$branch="master"};git checkout --orphan latest_branch;git add -A;git commit -am "init...";git branch -D $branch;git branch -m $branch;git push -f origin $branch;}
Function port{Param([Parameter(Mandatory=$true)][string]$from,[Parameter(Mandatory=$true)][string]$to);netsh interface portproxy reset;netsh interface portproxy add v4tov4 listenport=$to connectport=$from connectaddress=127.0.0.1}

Set-Alias -Name y -Value yarn
cls;
