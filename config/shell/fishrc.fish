status is-interactive; or return
set -g fish_greeting
set -l dotfiles_dir ~/.dotfiles
if test "$OS" = Windows_NT; and test -n "$DOTFILES_DIR"
    set dotfiles_dir (cygpath -u "$DOTFILES_DIR")
end

if test (uname) = Darwin
    if test -f $dotfiles_dir/.env.path
        while read -l p
            contains $p $PATH; or set -x PATH $PATH $p
        end < $dotfiles_dir/.env.path
    end
    if test -f ~/.env.path
        while read -l p
            contains $p $PATH; or set -x PATH $PATH $p
        end < ~/.env.path
    end
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
