if [[ $- == *i* ]]; then
	source $DOTFILES_DIR/config/shell/alias.sh
	eval "$(starship init bash)"
	eval "$(zoxide init bash)"
fi

export GOPATH=$(go env GOPATH)
export GOROOT=$(go env GOROOT)
export JAVA_HOME=$(mise where java)
