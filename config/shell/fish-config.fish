fish_config theme choose default
set fish_greeting
set fish_color_end normal
set fish_color_quote green
set fish_color_comment --dim
set fish_color_command magenta

if test "$OS" = Windows_NT
    dotsh fish (mise env --dotenv) | source
else
    /opt/homebrew/bin/brew shellenv fish | source
    mise activate fish | source
end

# Environment
test -f ~/.env; and dotsh fish "$(string collect < ~/.env)" | source
test -f ~/.path; and set -x PATH $PATH (cat ~/.path)
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
