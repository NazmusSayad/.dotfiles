function gac
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

function gp
    gac $argv
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

function git-reset
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
