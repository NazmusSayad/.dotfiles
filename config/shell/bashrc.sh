if [[ -f ~/.path ]]; then
	export PATH="$PATH:$(cat ~/.path)"
fi

if [[ "$OS" == "Windows_NT" ]]; then
	eval "$(dotsh bash "$(mise env --dotenv)")"
fi

if command -v uname >/dev/null 2>&1 && [[ "$(uname)" == "Darwin" ]]; then
	eval "$(brew shellenv bash)"
	eval "$(mise activate bash)"
fi

if [[ -f ~/.env ]]; then
	eval "$(dotsh bash "$(cat ~/.env)")"
fi

eval "$(direnv hook bash)"

eval "$(shaka bash)"
eval "$(zoxide init bash)"
eval "$(starship init bash)"

zoxide add "$PWD"

on_cd() {
	zoxide add "$PWD"
}

PROMPT_COMMAND="on_cd${PROMPT_COMMAND:+;$PROMPT_COMMAND}"
