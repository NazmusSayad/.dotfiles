if test "$OS" = Windows_NT
    dotsh fish (mise env --dotenv) | source
else
    brew shellenv fish | source
    mise activate fish | source
end

test -f ~/.env; and dotsh fish "$(cat ~/.env)" | source
test -f ~/.path; and set -x PATH $PATH (cat ~/.path)
direnv hook fish | source

shaka fish | source
zoxide init fish | source
starship init fish | source

zoxide add $PWD
function on_cd --on-variable PWD
    zoxide add $PWD
end
