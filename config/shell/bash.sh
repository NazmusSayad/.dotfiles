if [[ "$OS" == "Windows_NT" ]]; then
	eval "$(dotsh bash "$(mise env --dotenv)")"
else
	eval "$(/opt/homebrew/bin/brew shellenv bash)"
	eval "$(mise activate bash)"
fi

[[ -f ~/.env ]] && eval "$(dotsh bash "$(<~/.env)")"
[[ -f ~/.path ]] && export PATH="$PATH:$(paste -sd: ~/.path)"