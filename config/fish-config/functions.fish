function gbc
    set current (git branch --show-current)
    set branches (git branch --format="%(refname:short)" | grep -v $current)

    if test (count $branches) -gt 0
        set_color --dim
        echo -n "Branches to delete: "
        set_color normal

        set_color --bold red
        echo (string join ', ' $branches)
        set_color normal

        echo "Press [Enter] to confirm, or any other key to cancel: "
        read -n 1 confirm
        set read_status $status

        if test $read_status -ne 0
            set_color green
            echo "Cancelled branch deletion."
            set_color normal
            return 1
        end

        if test -z "$confirm"
            set_color --dim red
            git branch -d $branches
        else
            set_color green
            echo "Cancelled branch deletion."
        end

        set_color normal
    else
        echo "No other branches to delete"
    end
end


function gcommit
    if not git ls-files --error-unmatch -m --directory --no-empty-directory -o --exclude-standard ":/*" >/dev/null 2>&1
        echo "❌ No changes to commit."
        return 1
    end

    set msg "$argv"
    if test -z "$msg"
        echo "❗ No message provided, using git status as message."
        echo ""
        set msg (git status --porcelain)
    end

    git status --short; and echo
    git add -A >/dev/null; and git commit -m "$msg"
end


function gpush
    gcommit $argv
    if test $status -ne 0
        echo "❌ Commit failed or no changes to commit."
        return
    end

    if test -f package.json
        echo ""
        echo "❗ Running lint (if present)..."
        npm run lint --if-present

        if test $status -ne 0
            echo "❌ Linting failed. Fix the issues before pushing."
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


function fs-case
    fsutil file queryCaseSensitiveInfo .
end
