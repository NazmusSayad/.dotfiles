#!/usr/bin/env fish

set -g fish_greeting ''
set __dirname (realpath (dirname (status --current-filename)))

source $__dirname/aliases.fish
source $__dirname/functions.fish

if status is-interactive
    source $__dirname/theme.fish
end
