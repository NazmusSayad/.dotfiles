[[ $- != *i* ]] && return

if [[ "$(uname)" == "Darwin" ]]; then
	if [[ -f ~/.dotfiles/.env.path ]]; then
		while read -r p; do
			[[ ":$PATH:" != *":$p:"* ]] && export PATH="$PATH:$p"
		done <~/.dotfiles/.env.path
	fi
	if [[ -f ~/.env.path ]]; then
		while read -r p; do
			[[ ":$PATH:" != *":$p:"* ]] && export PATH="$PATH:$p"
		done <~/.env.path
	fi
fi

if [[ "$OS" == "Windows_NT" ]]; then
	eval "$(dotsh bash "$(mise env --dotenv)")"
fi

if [[ "$(uname)" == "Linux" ]]; then
	export PATH="$HOME/.local/bin:$PATH"
	eval "$("$HOME/.local/bin/mise" activate bash)"
fi

if [[ "$(uname)" == "Darwin" ]]; then
	eval "$(brew shellenv bash)"
	eval "$(mise activate bash)"
fi

[[ -f ~/.dotfiles/.env ]] && eval "$(dotsh bash "$(cat ~/.dotfiles/.env)")"
[[ -f ~/.env ]] && eval "$(dotsh bash "$(cat ~/.env)")"
eval "$(direnv hook bash)"

eval "$(shaka bash)"
eval "$(zoxide init bash)"
eval "$(starship init bash)"

zoxide add "$PWD" 2>/dev/null
on_cd() {
	zoxide add "$PWD" 2>/dev/null
}
PROMPT_COMMAND="on_cd${PROMPT_COMMAND:+;$PROMPT_COMMAND}"

if [[ "$OS" == "Windows_NT" ]]; then
	PROMPT_COMMAND=${PROMPT_COMMAND:+"$PROMPT_COMMAND; "}'printf "\e]9;9;%s\e\\" "`cygpath -w "$PWD" -C ANSI`"'
fi
