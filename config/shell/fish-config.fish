# Mise
if test "$OS" = Windows_NT
    dotsh fish (mise env -D) | source
else
    mise activate fish | source
end

# Environment
direnv hook fish | source

# Enhancements
shaka fish | source
zoxide init fish | source
starship init fish | source

# zoxide
zoxide add $PWD
function on_cd --on-variable PWD
    zoxide add $PWD
end
