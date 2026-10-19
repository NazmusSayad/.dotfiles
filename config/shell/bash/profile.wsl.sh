eval "$("$HOME/.local/bin/mise" env --shell bash)"
eval "$(dotsh init bash --home "${DOTFILES_DIR:-$HOME/.dotfiles}")"

. "$HOME/.bashrc"
