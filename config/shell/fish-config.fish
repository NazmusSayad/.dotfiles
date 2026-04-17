#!/usr/bin/env fish

fish_config theme choose default
set fish_color_end normal
set fish_color_quote green
set fish_color_comment --dim
set fish_color_command magenta

for line in (mise env --dotenv)
    set -l key (string split -m1 "=" $line)[1]
    set -l val (string split -m1 "=" $line)[2]
    set -gx $key $val
end

shaka fish | source
direnv hook fish | source

if status is-interactive
    zoxide init fish | source
    starship init fish | source
end

function on_cd --on-variable PWD
    zoxide add $PWD
end

function fish_greeting
    # fastfetch
end
