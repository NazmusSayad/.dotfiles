#!/usr/bin/env bash

source $DOTFILES_DIR/config/shell/alias.sh

eval "$(direnv hook bash)"
eval "$(mise env --dotenv)"

if [[ $- == *i* ]] && [[ "$OPENCODE_TERMINAL" != "1" ]]; then
	eval "$(starship init bash)"
	eval "$(zoxide init bash)"
fi
