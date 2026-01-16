if [[ -o interactive ]]; then
	source $DOTFILES_DIR/config/shell/alias.sh
	eval "$(starship init zsh)"
	eval "$(zoxide init zsh)"
fi

export GOBIN=$(go env GOBIN)
export GOROOT=$(go env GOROOT)
export JAVA_HOME=$(mise where java)
