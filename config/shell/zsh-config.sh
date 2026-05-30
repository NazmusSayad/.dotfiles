#!/usr/bin/env zsh

eval "$(dotsh zsh "$(mise env --dotenv)")"

if [[ $- == *i* ]]; then
	eval "$(shaka zsh)"
	eval "$(zoxide init zsh)"
	eval "$(starship init zsh)"
fi
