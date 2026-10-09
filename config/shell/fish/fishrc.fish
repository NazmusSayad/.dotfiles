status is-interactive; or return
set -g fish_greeting
set -l dotfiles_dir ~/.dotfiles
if test "$OS" = Windows_NT; and test -n "$DOTFILES_DIR"
    set dotfiles_dir (cygpath -u "$DOTFILES_DIR")
end

set -l path_file ''
if test "$OS" = Windows_NT
    set path_file $dotfiles_dir/.path.win
else if test (uname) = Darwin
    set path_file $dotfiles_dir/.path.mac
else if test -n "$WSL_DISTRO_NAME"; or string match -qi '*microsoft*' (uname -r)
    set path_file $dotfiles_dir/.path.wsl
end

if test -f "$path_file"
    while read -l p
        set p (string replace -r '\r$' '' -- "$p")
        test -z "$p"; and continue
        if test "$OS" = Windows_NT
            set p (cygpath -u "$p")
        end
        contains -- "$p" $PATH; or set -gx PATH $PATH "$p"
    end < "$path_file"
end

if test "$OS" = Windows_NT
    dotsh fish (mise env --dotenv) | source

    function __windows_terminal_report_cwd --on-event fish_prompt
        printf '\e]9;9;%s\e\\' (cygpath -w "$PWD" -C ANSI)
    end
end

if test (uname) = Linux
    set -gx PATH $HOME/.local/bin $PATH
    $HOME/.local/bin/mise activate fish | source
end

if command -q uname; and test (uname) = Darwin
    brew shellenv fish | source
    mise activate fish | source
end

test -f $dotfiles_dir/.env; and dotsh fish "$(cat $dotfiles_dir/.env)" | source
test -f $dotfiles_dir/.local/.env; and dotsh fish "$(cat $dotfiles_dir/.local/.env)" | source
direnv hook fish | source

shaka fish | source
zoxide init fish | source
starship init fish | source

zoxide add $PWD 2>/dev/null
function on_cd --on-variable PWD
    zoxide add $PWD 2>/dev/null
end
