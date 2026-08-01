[[ -f ~/.path ]] && export PATH="$PATH:$(paste -sd ':' ~/.path)"

eval "$(/opt/homebrew/bin/brew shellenv bash)"
eval "$(mise env --shell bash)"

[[ -f ~/.env ]] && eval "$(dotsh bash "$(cat ~/.env)")"
