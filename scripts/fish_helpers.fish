alias n="node --no-warnings"
alias nw="node --watch --no-warnings"

alias x="npm exec"
alias r="npm run"
alias rp="run-p --silent"
alias rs="run-s --silent"

alias ni="npm install"
alias nid="npm install --save-dev"
alias nu="npm uninstall"

alias gc="gh repo clone"

function gp
    set msg "$argv"

    if test -z "$msg"
        echo "No message provided, using git status as message:"
        set msg (git status --short --no-renames)
    end

    git status --short
    git add -A >/dev/null; and git commit -m "$msg"; and echo; and git push --quiet
end


gp