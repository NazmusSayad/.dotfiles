if [[ $- == *i* ]]; then
	source $DOTFILES_DIR/config/shell/alias.sh
	eval "$(starship init bash)"
	eval "$(zoxide init bash)"
fi
