export GOBIN=$(go env GOBIN)
export GOROOT=$(go env GOROOT)
export JAVA_HOME=$(mise where java)

if [[ $- == *i* ]]; then
	source $DOTFILES_DIR/config/shell/alias.sh
	eval "$(starship init bash)"
	eval "$(zoxide init bash)"
fi
