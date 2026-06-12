if test "$OS" = Windows_NT
    dotsh fish (mise env -D) | source
else
    /opt/homebrew/bin/brew shellenv fish | source
    mise activate fish | source
end

direnv hook fish | source

function fish_greeting
end

if status is-interactive
    fish_config theme choose default
    set fish_color_end normal
    set fish_color_quote green
    set fish_color_comment --dim
    set fish_color_command magenta

    shaka fish | source
    zoxide init fish | source
    starship init fish | source

    zoxide add $PWD
    function on_cd --on-variable PWD
        zoxide add $PWD
    end
end
