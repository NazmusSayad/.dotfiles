if [[ -x /opt/homebrew/bin/brew ]]; then
	eval "$(/opt/homebrew/bin/brew shellenv bash)"
fi

if command -v mise >/dev/null 2>&1; then
	eval "$(mise env bash)"
fi

[[ -f ~/.env ]] && eval "$(dotsh bash "$(<~/.env)")"
[[ -f ~/.path ]] && export PATH="$PATH:$(paste -sd: ':' ~/.path)"
