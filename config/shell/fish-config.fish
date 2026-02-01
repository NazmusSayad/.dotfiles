#!/usr/bin/env fish

fish_config theme choose default
set fish_color_end normal
set fish_color_quote green
set fish_color_comment --dim
set fish_color_command magenta

set -g GOBIN (go env GOBIN)
set -g GOROOT (go env GOROOT)
set -g JAVA_HOME (mise where java)

direnv hook fish | source

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
