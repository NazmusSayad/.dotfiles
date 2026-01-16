#!/usr/bin/env fish

set -g fish_color_command magenta

source $DOTFILES_DIR/config/shell/alias.sh

if status is-interactive
    starship init fish | source
    zoxide init fish | source
end

function fish_greeting
    if not set -q TERM_PROGRAM
        # fastfetch
    end
end

set -g GOPATH (go env GOPATH)
set -g GOROOT (go env GOROOT)
set -g JAVA_HOME (mise where java)
