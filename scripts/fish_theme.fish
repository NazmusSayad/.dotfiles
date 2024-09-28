function __fish_shell_path
    set dir (pwd)
    if test "$OSTYPE" = cygwin -o "$OSTYPE" = msys
        set win_dir $dir
        set first_part (string sub -s 2 -l 1 $win_dir)
        set first_part_uppercase (string upper $first_part)
        set second_part (string sub -s 4 $win_dir)
        set dir "$first_part_uppercase: $second_part"
    end
    echo $dir
end

function __fish_git_branch
    if not git rev-parse --is-inside-work-tree >/dev/null 2>&1
        return
    end

    set branch (git symbolic-ref --short HEAD 2>/dev/null; or echo (git rev-parse --short HEAD 2>/dev/null))
    set output "$branch"
    echo -n "$branch"
end

function __fish_git_changes
    if git ls-files --error-unmatch -m --directory --no-empty-directory -o --exclude-standard ":/*" >/dev/null 2>&1
        set output "$output ✗"
    end
    set output "$output"
    echo -n "$output"
end

function fish_prompt
    set_color magenta
    echo -n '⌠ '
    set_color blue
    echo -n (__fish_shell_path)
    set_color normal

    set branch_output (__fish_git_branch)
    if test -n "$branch_output"
        echo -n ' ('
        set_color yellow
        echo -n "$branch_output"
        set_color normal

        set changes_output (__fish_git_changes)
        if test -n $changes_output
            set_color red
            echo -n $changes_output
            set_color normal
        end

        echo -n ')'
    end

    echo
    set_color magenta
    echo -n '$ '
    set_color normal
end
