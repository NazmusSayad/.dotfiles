#!/usr/bin/env fish

set -g fish_color_command magenta

set -g GOBIN (go env GOBIN)
set -g GOROOT (go env GOROOT)
set -g JAVA_HOME (mise where java)

function fish_greeting
    if not set -q TERM_PROGRAM
        # fastfetch
    end
end

if status is-interactive
    source $DOTFILES_DIR/config/shell/alias.sh
    starship init fish | source
    zoxide init fish | source
end
