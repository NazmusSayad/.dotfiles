#!/usr/bin/env bash

eval "$(shell-alias sh)"
eval "$(direnv hook bash)"
eval "$(mise env --dotenv)"

if [[ $- == *i* ]]; then
	eval "$(starship init bash)"
	eval "$(zoxide init bash)"
fi
