if [[ -f ~/.env ]]; then
	eval "$(dotsh bash "$(<~/.env)")"
fi

if [[ -f ~/.path ]]; then
	export PATH="$PATH:$(paste -s -d ':' ~/.path)"
fi

if [[ "$OS" == "Windows_NT" ]]; then
	eval "$(dotsh bash "$(mise env --dotenv)")"
else
	if [[ -x /opt/homebrew/bin/brew ]]; then
		eval "$(/opt/homebrew/bin/brew shellenv bash)"
	fi

	if command -v mise >/dev/null 2>&1; then
		eval "$(mise env bash)"
	fi
fi
