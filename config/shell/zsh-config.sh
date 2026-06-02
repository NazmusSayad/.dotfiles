eval "$(dotsh zsh "$(mise env --dotenv)")"
eval "$(direnv hook zsh)"

if [[ $- == *i* ]]; then
	eval "$(shaka zsh)"
	eval "$(zoxide init zsh)"
	eval "$(starship init zsh)"
fi
