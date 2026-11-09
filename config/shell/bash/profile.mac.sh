eval "$(/opt/homebrew/bin/brew shellenv bash)"
eval "$("$HOME/.local/bin/mise" env --shell bash)"
eval "$(dotsh init bash --home "${DOTSYS_DIR:-$HOME/.dotsys}")"

. "$HOME/.bashrc"
