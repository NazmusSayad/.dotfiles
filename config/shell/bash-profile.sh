export RUST_BACKTRACE=1
export NODE_NO_WARNINGS=1

eval "$(/opt/homebrew/bin/brew shellenv bash)"
eval "$(mise env --shell bash)"

if [[ -f ~/.env ]]; then
	eval "$(dotsh bash "$(cat ~/.env)")"
fi

if [[ -f ~/.path ]]; then
	export PATH="$(paste -s -d ':' ~/.path):$PATH"
fi
