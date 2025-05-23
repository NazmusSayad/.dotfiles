function __fish_shell_path
    set dir (pwd)
    if test "$OS" = Windows_NT
        set win_dir $dir
        set first_part (string sub -s 2 -l 1 $win_dir)
        set first_part_uppercase (string upper $first_part)
        set second_part (string sub -s 4 $win_dir)
        set dir "$first_part_uppercase: $second_part"
    end
    echo -n $dir
end

function __fish_git_branch
    if not git rev-parse --is-inside-work-tree >/dev/null 2>&1
        return
    end

    echo -n (git symbolic-ref --short HEAD 2>/dev/null; or echo (git rev-parse --short HEAD 2>/dev/null))
end

function __fish_git_changes
    if not git diff --quiet --cached || not git diff --quiet || git ls-files --others --exclude-standard | grep -q "."
        echo -n "✘"
    else if git log --branches --not --remotes | grep -q "."
        echo -n "✚"
    end
end

function fish_prompt
    set_color magenta
    echo -n (__fish_shell_path)
    set_color normal

    set branch_output (__fish_git_branch)
    if test -n "$branch_output"
        echo -n " ("
        set_color red
        echo -n "$branch_output"
        set_color normal

        set changes_output (__fish_git_changes)
        if test -n "$changes_output"
            set_color yellow
            echo -n " $changes_output"
            set_color normal
        end

        echo -n ")"
    end

    set_color magenta
    echo -en "\n⎩ "
    set_color normal
end
