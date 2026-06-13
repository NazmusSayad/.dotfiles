if [[ "$OS" == "Windows_NT" ]]; then
	eval "$(dotsh bash "$(mise env --dotenv)")"
else
	eval "$(/opt/homebrew/bin/brew shellenv bash)"
	eval "$(mise activate bash)"
	export PATH="$PATH:~/.dotfiles/.build/bin"
fi

[[ -f ~/.env ]] && eval "$(dotsh bash "$(<~/.env)")"
eval "$(direnv hook bash)"

if [[ $- == *i* ]]; then
	eval "$(shaka bash)"
	eval "$(zoxide init bash)"
	eval "$(starship init bash)"
fi
