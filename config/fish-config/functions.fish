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


function gcommit
    if not git ls-files --error-unmatch -m --directory --no-empty-directory -o --exclude-standard ":/*" >/dev/null 2>&1
        echo "❌ No changes to commit"
        return 1
    end

    set msg "$argv"
    if test -z "$msg"
        echo "❗ No message provided, using git status as message"
        echo ""
        set msg (git status --porcelain)
    end

    git status --short; and echo
    git add -A >/dev/null; and git commit -m "$msg"
end


function gpush
    gcommit $argv
    if test $status -ne 0
        echo "❌ Commit failed or no changes to commit"
        return
    end

    if test -f package.json
        echo ""
        echo "❗ Running lint (if present)..."
        npm run lint --if-present

        if test $status -ne 0
            echo "❌ Linting failed. Fix the issues before pushing"
            return
        end
    end

    echo ""
    git push
end


function greset
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


function gpm
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


function gpr
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


function gr
    set_color red
    echo "Reset and clean?"

    set_color normal
    set_color --dim
    echo -n "Press [Enter] to confirm, or any other key to cancel: "
    set_color normal

    read -n 1 -P "" --function confirm

    if test $status -eq 0 -a -z "$confirm"
        echo "✅ Done."
    else
        echo "❌ Cancelled."
    end
end


function gpg-unlock
    for pid in (ps aux | grep gpg | grep -v grep | awk '{print $1}')
        echo "Found GPG process with PID: $pid"
        kill -9 $pid
    end
    rm -f ~/.gnupg/*.lock
end


function fscase
    fsutil.exe file setCaseSensitiveInfo . enable recursive
end
