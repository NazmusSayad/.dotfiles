function __fish_shell_path
    set -l dir (pwd)

    if string match -q "$HOME/Desktop*" $dir
        if string match -q "$HOME/Desktop" $dir
            set dir "🖥️ DESKTOP"
        else
            set dir (string replace "$HOME/Desktop/" "🖥️ " $dir)
        end

    else if string match -q "$HOME/Documents*" $dir
        if string match -q "$HOME/Documents" $dir
            set dir "📄 DOCUMENTS"
        else
            set dir (string replace "$HOME/Documents/" "📄 " $dir)
        end

    else if string match -q "$HOME/Downloads*" $dir
        if string match -q "$HOME/Downloads" $dir
            set dir "📥 DOWNLOADS"
        else
            set dir (string replace "$HOME/Downloads/" "📥 " $dir)
        end

    else if string match -q "$HOME/Pictures*" $dir
        if string match -q "$HOME/Pictures" $dir
            set dir "🖼️ PICTURES"
        else
            set dir (string replace "$HOME/Pictures/" "🖼️ " $dir)
        end

    else if string match -q "$HOME/Videos*" $dir
        if string match -q "$HOME/Videos" $dir
            set dir "🎥 VIDEOS"
        else
            set dir (string replace "$HOME/Videos/" "🎥 " $dir)
        end

    else if string match -q "$HOME/Music*" $dir
        if string match -q "$HOME/Music" $dir
            set dir "🎵 MUSIC"
        else
            set dir (string replace "$HOME/Music/" "🎵 " $dir)
        end

    else if test "$OS" = Windows_NT
        set win_dir $dir
        set first_part (string sub -s 2 -l 1 $win_dir)
        set first_part_uppercase (string upper $first_part)
        set second_part (string sub -s 4 $win_dir)
        set dir "$first_part_uppercase: $second_part"
    end

    echo -n $dir
end

function __fish_git_changes
    if not git diff --quiet --cached || not git diff --quiet || git ls-files --others --exclude-standard | grep -q "."
        echo -n "*"
    else if git log --branches --not --remotes | grep -q "."
        echo -n "+"
    end
end

function fish_prompt
    set_color magenta
    echo -n (__fish_shell_path)
    set_color normal

    set branch_output (git branch --show-current 2>/dev/null)
    if test -n "$branch_output"
        set_color --dim
        echo -n " ("
        set_color normal

        set_color red
        echo -n "$branch_output"
        set_color normal

        set changes_output (__fish_git_changes)
        if test -n "$changes_output"
            set_color yellow --bold
            echo -n "$changes_output"
            set_color normal
        end

        set_color --dim
        echo -n ")"
        set_color normal

    end

    set_color magenta
    echo -en "\n⎩ "
    set_color normal
end
