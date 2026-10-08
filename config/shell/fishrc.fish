status is-interactive; or return
if test (uname) = Darwin
    if test -f ~/.dotfiles/.env.path
        while read -l p
            contains $p $PATH; or set -x PATH $PATH $p
        end < ~/.dotfiles/.env.path
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

if command -q uname; and test (uname) = Darwin
    brew shellenv fish | source
    mise activate fish | source
end

test -f ~/.dotfiles/.env; and dotsh fish "$(cat ~/.dotfiles/.env)" | source
test -f ~/.env; and dotsh fish "$(cat ~/.env)" | source
direnv hook fish | source

shaka fish | source
zoxide init fish | source
starship init fish | source

zoxide add $PWD 2>/dev/null
function on_cd --on-variable PWD
    zoxide add $PWD 2>/dev/null
end
