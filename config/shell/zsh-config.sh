#!/usr/bin/env zsh

source $DOTFILES_DIR/config/shell/alias.sh
source $DOTFILES_DIR/config/shell/mise-env.sh

eval "$(direnv hook zsh)"

if [[ -o interactive ]]; then
	eval "$(starship init zsh)"
	eval "$(zoxide init zsh)"
fi
