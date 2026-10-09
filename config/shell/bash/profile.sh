if [[ "$(uname)" == "Darwin" ]]; then
	if [[ -f ~/.dotfiles/.path.mac ]]; then
		while IFS= read -r p || [[ -n "$p" ]]; do
			p="${p%$'\r'}"
			[[ -z "$p" ]] && continue
			[[ ":$PATH:" != *":$p:"* ]] && export PATH="$PATH:$p"
		done <~/.dotfiles/.path.mac
	fi
fi

eval "$(/opt/homebrew/bin/brew shellenv bash)"
eval "$(mise env --shell bash)"

[[ -f ~/.dotfiles/.env ]] && eval "$(dotsh bash "$(cat ~/.dotfiles/.env)")"
[[ -f ~/.dotfiles/.local/.env ]] && eval "$(dotsh bash "$(cat ~/.dotfiles/.local/.env)")"
