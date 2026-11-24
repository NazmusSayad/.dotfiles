[[ ":$PATH:" != *":$HOME/.local/bin:"* ]] && export PATH="$HOME/.local/bin:$PATH"

eval "$("$HOME/.local/bin/mise" env --shell bash)"
eval "$(dotsh init bash --home "${DOTFILES_DIR:-$HOME/.dotfiles}/.dotsys")"

. "$HOME/.bashrc"
