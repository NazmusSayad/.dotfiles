#!/usr/bin/env bash

source $DOTFILES_DIR/config/shell/alias.sh

export GOBIN=$(go env GOBIN)
export GOROOT=$(go env GOROOT)
export JAVA_HOME=$(mise where java)

eval "$(direnv hook bash)"

if [[ $- == *i* ]]; then
	eval "$(starship init bash)"
	eval "$(zoxide init bash)"
fi
