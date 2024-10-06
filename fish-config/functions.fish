function gp
    if not git ls-files --error-unmatch -m --directory --no-empty-directory -o --exclude-standard ":/*" >/dev/null 2>&1
        echo "No changes to commit."
        return
    end

    set msg "$argv"
    if test -z "$msg"
        echo "No message provided, using git status as message."
        set msg (git status --porcelain)
    end

    git status --short; and echo
    git add -A >/dev/null; and git commit -m "$msg"; and echo; and git push --quiet
end

function gr
    set hash $argv[1]
    if test -z "$hash"
        echo "You must need to give a <commit_id>"
        return
    end

    git checkout $hash .; and gp "$hash restored; $argv[2..-1]"
end

function ghr
    set branch $argv[1]
    if test -z "$branch"
        set branch master
    end

    git checkout --orphan latest_branch
    git add -A
    git commit -am "initial commit"
    git branch -D $branch
    git branch -m $branch
    git push -f origin $branch
end


function gpull
    set src $argv[1]

    if test -z "$src"
        echo "Error: Please provide the source branch."
        return 1
    end

    if test ! -d ".git"
        echo "Error: Not a Git repository."
        return 1
    end

    set target (git rev-parse --abbrev-ref HEAD)
    git pull origin "$src:$target"
    return $status
end
