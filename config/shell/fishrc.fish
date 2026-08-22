status is-interactive; or return
if test -f ~/.path
    while read -l p
        contains $p $PATH; or set -x PATH $PATH $p
    end < ~/.path
end

if test "$OS" = Windows_NT
    dotsh fish (mise env --dotenv) | source
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

zoxide add $PWD
function on_cd --on-variable PWD
    zoxide add $PWD
end
