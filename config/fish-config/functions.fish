function gbc # git branch clean
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
            git branch -d $branches
        else
            set_color green
            echo "Cancelled branch deletion"
        end

        set_color normal
    else
        echo "No other branches to delete"
    end
end

function greset # git complete reset
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

function gpm # git pull merge
    if test (count $argv) -ne 1
        echo "Usage: gpm <branch>" >&2
        return 1
    end

    set target_branch $argv[1]
    set current_branch (git branch --show-current)

    echo -n "Pulling changes from "
    set_color red
    echo -n "$target_branch"
    set_color normal
    echo -n " into "
    set_color blue
    echo -n "$current_branch"
    set_color normal
    echo " (merge)..."
    git pull origin $target_branch --no-rebase
end

function gpr # git pull rebase
    if test (count $argv) -ne 1
        echo "Usage: gpr <branch>" >&2
        return 1
    end

    set target_branch $argv[1]
    set current_branch (git branch --show-current)

    echo -n "Pulling changes from "
    set_color red
    echo -n "$target_branch"
    set_color normal
    echo -n " into "
    set_color blue
    echo -n "$current_branch"
    set_color normal
    echo " (rebase)..."
    git pull origin $target_branch --rebase
end

function gr # git restore
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
        echo "❌ Cancelled."
    end
end

function gpg-unlock # gpg unlock
    for pid in (ps aux | grep gpg | grep -v grep | awk '{print $1}')
        echo "Found GPG process with PID: $pid"
        kill -9 $pid
    end

    for lf in ~/.gnupg/*.lock
        rm -f $lf
    end
end

function fs-case # fsutil.exe file setCaseSensitiveInfo . enable recursive
    fsutil.exe file setCaseSensitiveInfo . enable recursive
end
