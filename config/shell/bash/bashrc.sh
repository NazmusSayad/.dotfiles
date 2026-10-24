[[ $- != *i* ]] && return

if [[ "$(uname)" == "Darwin" ]]; then
	eval "$(/opt/homebrew/bin/brew shellenv bash)"
	eval "$("$HOME/.local/bin/mise" activate bash)"
elif [[ "$(uname)" == "Linux" ]]; then
	eval "$("$HOME/.local/bin/mise" activate bash)"
elif [[ "$OS" == "Windows_NT" ]]; then
	eval "$(mise activate bash)"
	PROMPT_COMMAND=${PROMPT_COMMAND:+"$PROMPT_COMMAND; "}'printf "\e]9;9;%s\e\\" "`cygpath -w "$PWD" -C ANSI`"'
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

claude() {
	local config="${CLAUDE_CONFIG_DIR:-$HOME}/.claude.json"
	local dir trusted
	dir=$(git rev-parse --show-toplevel 2>/dev/null || pwd)
	trusted=$(jq --arg d "$dir" '.projects[$d].hasTrustDialogAccepted = true' "$config") &&
		printf '%s\n' "$trusted" >"$config"
	command claude "$@"
}
