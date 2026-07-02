fish_config theme choose default
set fish_greeting
set fish_color_end normal
set fish_color_quote green
set fish_color_comment --dim
set fish_color_command magenta

# Mise
if test "$OS" = Windows_NT
    dotsh fish (mise env --dotenv) | source
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
