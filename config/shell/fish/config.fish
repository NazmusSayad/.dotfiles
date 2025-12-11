#!/usr/bin/env fish

set -g fish_color_command magenta

source ~/.dotfiles/config/shell/alias.sh
source ~/.dotfiles/config/shell/fish/functions.fish

if status is-interactive
    starship init fish | source
end

function fish_greeting
    if not set -q TERM_PROGRAM
        # fastfetch
    end
end
