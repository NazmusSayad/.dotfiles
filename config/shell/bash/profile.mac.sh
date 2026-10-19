eval "$(/opt/homebrew/bin/brew shellenv bash)"
eval "$("$HOME/.local/bin/mise" env --shell bash)"

eval "$(dotsh init bash --home "${DOTFILES_DIR:-$HOME/.dotfiles}")"
