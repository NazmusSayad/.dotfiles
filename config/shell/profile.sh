if [[ "$(uname)" == "Darwin" ]]; then
[[ -f ~/.dotfiles/.env.path ]] && export PATH="$PATH:$(paste -sd ':' ~/.dotfiles/.env.path)"
[[ -f ~/.env.path ]] && export PATH="$PATH:$(paste -sd ':' ~/.env.path)"
fi

eval "$(/opt/homebrew/bin/brew shellenv bash)"
eval "$(mise env --shell bash)"

[[ -f ~/.dotfiles/.env ]] && eval "$(dotsh bash "$(cat ~/.dotfiles/.env)")"
[[ -f ~/.dotfiles/.local/.env ]] && eval "$(dotsh bash "$(cat ~/.dotfiles/.local/.env)")"
