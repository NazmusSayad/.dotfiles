[[ $- != *i* ]] && return

dotfiles_dir="$HOME/.dotfiles"
if [[ "$OS" == "Windows_NT" && -n "$DOTFILES_DIR" ]]; then
	dotfiles_dir="$(cygpath -u "$DOTFILES_DIR")"
fi

path_file=""
if [[ "$OS" == "Windows_NT" ]]; then
	path_file="$dotfiles_dir/.path.win"
elif [[ "$(uname)" == "Darwin" ]]; then
	path_file="$dotfiles_dir/.path.mac"
elif [[ -n "$WSL_DISTRO_NAME" || "$(uname -r)" == *[Mm]icrosoft* ]]; then
	path_file="$dotfiles_dir/.path.wsl"
fi

if [[ -f "$path_file" ]]; then
	while IFS= read -r p || [[ -n "$p" ]]; do
		p="${p%$'\r'}"
		[[ -z "$p" ]] && continue
		[[ "$OS" == "Windows_NT" ]] && p="$(cygpath -u "$p")"
		[[ ":$PATH:" != *":$p:"* ]] && export PATH="$PATH:$p"
	done <"$path_file"
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

[[ -f "$dotfiles_dir/.env" ]] && eval "$(dotsh bash "$(cat "$dotfiles_dir/.env")")"
[[ -f "$dotfiles_dir/.local/.env" ]] && eval "$(dotsh bash "$(cat "$dotfiles_dir/.local/.env")")"
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
