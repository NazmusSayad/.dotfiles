if test -f ~/.path
    set -x PATH $PATH (cat ~/.path)
end

if test "$OS" = Windows_NT
    dotsh fish (mise env --dotenv) | source
end

if command -q uname; and test (uname) = Darwin
    brew shellenv fish | source
    mise activate fish | source
end

if test -f ~/.env
    dotsh fish "$(cat ~/.env)" | source
end

direnv hook fish | source

shaka fish | source
zoxide init fish | source
starship init fish | source

zoxide add $PWD
function on_cd --on-variable PWD
    zoxide add $PWD
end
