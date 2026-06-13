if [[ "$OS" == "Windows_NT" ]]; then
	eval "$(dotsh zsh "$(mise env --dotenv)")"
else
	eval "$(/opt/homebrew/bin/brew shellenv zsh)"
	eval "$(mise activate zsh)"
fi

[[ -f ~/.env ]] && eval "$(dotsh zsh "$(<~/.env)")"
eval "$(direnv hook zsh)"

if [[ $- == *i* ]]; then
	eval "$(shaka zsh)"
	eval "$(zoxide init zsh)"
	eval "$(starship init zsh)"
fi
