eval "$(/opt/homebrew/bin/brew shellenv bash)"
eval "$(mise activate bash)"
eval "$(direnv hook bash)"

if [[ $- == *i* ]]; then
	eval "$(shaka bash)"
	eval "$(zoxide init bash)"
	eval "$(starship init bash)"
fi
