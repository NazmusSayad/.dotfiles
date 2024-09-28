set -g fish_greeting ''

source ~/.config/fish/fish_aliases.fish
source ~/.config/fish/fish_helpers.fish

if status is-interactive
    source ~/.config/fish/fish_theme.fish
end
