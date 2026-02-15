#!/usr/bin/env bash

source $DOTFILES_DIR/config/shell/alias.sh
source $DOTFILES_DIR/config/shell/mise-env.sh

eval "$(direnv hook bash)"

if [[ $- == *i* ]]; then
	eval "$(starship init bash)"
	eval "$(zoxide init bash)"
fi
