#!/usr/bin/env bash

eval "$(mise env --dotenv)"

eval "$(shaka bash)"
eval "$(direnv hook bash)"

if [[ $- == *i* ]]; then
	eval "$(zoxide init bash)"
	eval "$(starship init bash)"
fi
