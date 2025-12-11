#!/usr/bin/env fish

set -g fish_color_command magenta

if status is-interactive
    starship init fish | source
end

function fish_greeting
    if not set -q TERM_PROGRAM
        fastfetch
    end
end

source ~/.dotfiles/config/fish/aliases.fish
source ~/.dotfiles/config/fish/functions.fish
