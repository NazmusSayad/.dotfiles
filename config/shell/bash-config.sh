#!/usr/bin/env bash

eval "$(shaka bash)"
eval "$(direnv hook bash)"

eval "$(mise env --dotenv)"

if [[ $- == *i* ]]; then
	eval "$(zoxide init bash)"
	eval "$(starship init bash)"
fi
