eval "$(/opt/homebrew/bin/brew shellenv zsh)"
eval "$(mise activate zsh)"
eval "$(direnv hook zsh)"

if [[ $- == *i* ]]; then
	eval "$(shaka zsh)"
	eval "$(zoxide init zsh)"
	eval "$(starship init zsh)"
fi
