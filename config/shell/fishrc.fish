status is-interactive; or return

if test "$OS" = Windows_NT
    dotsh fish (mise env --dotenv) | source
end

if command -q uname; and test (uname) = Darwin
    brew shellenv fish | source
    mise activate fish | source
end

test -f ~/.path; and set -x PATH $PATH (cat ~/.path)
test -f ~/.env; and dotsh fish "$(cat ~/.env)" | source
direnv hook fish | source

shaka fish | source
zoxide init fish | source
starship init fish | source

zoxide add $PWD
function on_cd --on-variable PWD
    zoxide add $PWD
end
