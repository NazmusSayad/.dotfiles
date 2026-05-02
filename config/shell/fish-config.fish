#!/usr/bin/env fish

function fish_greeting
end

for line in (mise env --dotenv)
    set -l key (string split -m1 "=" $line)[1]
    set -l val (string split -m1 "=" $line)[2]
    set -gx $key $val
end

if status is-interactive
    # Setup fish theme
    fish_config theme choose default
    set fish_color_end normal
    set fish_color_quote green
    set fish_color_comment --dim
    set fish_color_command magenta

    shaka fish | source
    zoxide init fish | source
    starship init fish | source
end

function on_cd --on-variable PWD
    zoxide add $PWD
end
