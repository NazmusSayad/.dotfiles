[[ $- != *i* ]] && return

if [[ "$(uname)" == "Darwin" ]]; then
	eval "$(/opt/homebrew/bin/brew shellenv bash)"
	eval "$("$HOME/.local/bin/mise" activate bash)"
elif [[ "$(uname)" == "Linux" ]]; then
	eval "$("$HOME/.local/bin/mise" activate bash)"
elif [[ "$OS" == "Windows_NT" ]]; then
	eval "$(dotsh bash "$(mise env --dotenv)")"
fi

eval "$(dotsh init bash --home "${DOTFILES_DIR:-$HOME/.dotfiles}")"
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
