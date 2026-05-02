#!/usr/bin/env fish

dotsh fish (mise env -D) | source

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

    function on_cd --on-variable PWD
        zoxide add $PWD
    end
end
