# Git Branch Cleanup
function gbc
    set current (git branch --show-current)
    set branches (git branch --format="%(refname:short)" | grep -v $current)

    if test (count $branches) -gt 0
        echo -n "Branches to delete: "

        set_color --bold red
        echo (string join ', ' $branches)
        set_color normal

        set_color normal
        set_color --dim
        echo -n "Press [Enter] to confirm, or any other key to cancel: "
        set_color normal

        read -n 1 -P "" --function confirm

        if test $status -eq 0 -a -z "$confirm"
            set_color --dim red
            git prune --progress
            git branch -D $branches
        else
            set_color green
            echo "Cancelled branch deletion"
        end

        set_color normal
    else
        echo "No other branches to delete"
    end
end

# Git Reset
function greset
    set_color red
    echo "This will reset the entire repository to the latest remote branch."
    set_color normal

    echo "Write 'yes' and press [Enter] to confirm."
    read -P "> " confirm

    if test "$confirm" != yes
        set_color green
        echo "Reset aborted"
        set_color normal
        return 0
    end

    git fetch --all

    set --local remote_url (git remote get-url origin)
    set --local current_branch (git branch --show-current)

    echo -en "> Branch: "
    set_color yellow
    echo -e "$current_branch"
    set_color normal

    echo -en "> Remote: "
    set_color blue
    echo -e "$remote_url"
    set_color normal

    set --local remote_branches (git branch -r --format="%(refname:short)" | string split '\n')
    for rb in $remote_branches
        set --local match (string match -r '.+/.+' $rb)
        if test -z "$match"
            continue
        end

        set --local rb (string replace -r '^[^/]*/' '' $rb)
        if test "$rb" = "$current_branch"
            continue
        end

        echo -en "> Deleting remote branch: "
        set_color red
        echo -e "$rb"
        set_color normal

        git push origin --delete $rb
    end

    set_color red
    echo '> Deleting git folder...'
    set_color normal

    rm -rf .git

    git init --initial-branch=$current_branch
    git remote add origin $remote_url

    git add .
    git commit -m "Initial commit"
    git push --force --set-upstream origin $current_branch

end

# Git Restore
function gr
    set_color red
    echo "Restore and clean?"

    set_color normal
    set_color --dim
    echo -n "Press [Enter] to confirm, or any other key to cancel: "
    set_color normal

    read -n 1 -P "" --function confirm

    if test $status -eq 0 -a -z "$confirm"
        git restore .
        git clean -fd
    else
        set_color red
        echo "❌ Aborted."
        set_color normal
        return 0
    end
end

# Git Pull (Default)
function gp
    set current_branch (git branch --show-current)

    if test (count $argv) -eq 0
        set_color normal --dim
        echo "No branch specified, using current branch"
        set_color normal
        set target_branch $current_branch
    else if test (count $argv) -eq 1
        set target_branch $argv[1]
    else
        echo "Usage: gp [branch]" >&2
        return 1
    end

    set_color normal --dim
    echo -n "Pulling changes from "
    set_color blue --dim
    echo -n "$target_branch"
    set_color normal --dim
    echo -n " into "
    set_color red --dim
    echo -n "$current_branch"
    set_color normal --dim
    echo " (default)"
    set_color normal

    git prune --progress
    git pull origin $target_branch --progress
end

# Git Pull (Rebase)
function gpr
    set current_branch (git branch --show-current)

    if test (count $argv) -eq 0
        set_color normal --dim
        echo "No branch specified, using current branch"
        set_color normal
        set target_branch $current_branch
    else if test (count $argv) -eq 1
        set target_branch $argv[1]
    else
        echo "Usage: gp [branch]" >&2
        return 1
    end

    set_color normal --dim
    echo -n "Pulling changes from "
    set_color blue --dim
    echo -n "$target_branch"
    set_color normal --dim
    echo -n " into "
    set_color red --dim
    echo -n "$current_branch"
    set_color normal --dim
    echo " (rebase)"
    set_color normal

    git prune --progress
    git pull origin $target_branch --progress --rebase
end

# GPG Unlock
function gpg-unlock
    for pid in (ps aux | grep gpg | grep -v grep | awk '{print $1}')
        echo "Found GPG process with PID: $pid"
        kill -9 $pid
    end

    for pid in (ps aux | grep keyboxd | grep -v grep | awk '{print $1}')
        echo "Found keyboxd process with PID: $pid"
        kill -9 $pid
    end

    for lf in ~/.gnupg/*.lock
        rm -f $lf
    end
end

# File System Case
function fs-case
    fsutil.exe file setCaseSensitiveInfo . enable recursive
end
