#!/usr/bin/env fish

fish_config theme choose default
set fish_color_end normal
set fish_color_quote green
set fish_color_comment --dim
set fish_color_command magenta

function fish_greeting
    # fastfetch
end

source $DOTFILES_DIR/config/shell/alias.sh

for line in (mise env --dotenv)
    set -l key (string split -m1 "=" $line)[1]
    set -l val (string split -m1 "=" $line)[2]
    set -gx $key $val
end

direnv hook fish | source

if status is-interactive
    starship init fish | source
    zoxide init fish | source
end
