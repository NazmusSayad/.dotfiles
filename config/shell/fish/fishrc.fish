status is-interactive; or return
set -g fish_greeting

if command -q uname; and test (uname) = Darwin
    /opt/homebrew/bin/brew shellenv fish | source
    "$HOME/.local/bin/mise" activate fish | source
else if test (uname) = Linux
    "$HOME/.local/bin/mise" activate fish | source
else if test "$OS" = Windows_NT
    dotsh fish (mise env --dotenv) | source

    function __windows_terminal_report_cwd --on-event fish_prompt
        printf '\e]9;9;%s\e\\' (cygpath -w "$PWD" -C ANSI)
    end
end

set -l dotfiles_dir "$HOME/.dotfiles"
if test -n "$DOTFILES_DIR"
    set dotfiles_dir "$DOTFILES_DIR"
end
dotsh init fish --home "$dotfiles_dir" | source
direnv hook fish | source

shaka fish | source
zoxide init fish | source
starship init fish | source

zoxide add $PWD 2>/dev/null
function on_cd --on-variable PWD
    zoxide add $PWD 2>/dev/null
end
