#!/usr/bin/env bash

eval "$(mise env --dotenv)"

if [[ $- == *i* ]]; then
	eval "$(shaka bash)"
	eval "$(zoxide init bash)"
	eval "$(starship init bash)"
fi
