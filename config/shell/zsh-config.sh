#!/usr/bin/env zsh

source $DOTFILES_DIR/config/shell/alias.sh

eval "$(direnv hook zsh)"
eval "$(mise env --dotenv)"

if [[ -o interactive ]]; then
	eval "$(starship init zsh)"
	eval "$(zoxide init zsh)"
fi
