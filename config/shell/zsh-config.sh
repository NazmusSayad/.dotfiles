#!/usr/bin/env zsh

source $DOTFILES_DIR/config/shell/alias.sh

export GOBIN=$(go env GOBIN)
export GOROOT=$(go env GOROOT)
export JAVA_HOME=$(mise where java)

eval "$(direnv hook zsh)"

if [[ -o interactive ]]; then
	eval "$(starship init zsh)"
	eval "$(zoxide init zsh)"
fi
