eval "$(/opt/homebrew/bin/brew shellenv bash)"
eval "$("$HOME/.local/bin/mise" env --shell bash)"
eval "$(dotsh init bash --home "${DOTSYS_DIR:-$HOME/.dotsys}")"

. "$HOME/.bashrc"

if [ -n "$SSH_CONNECTION" ] && [ -t 0 ]; then
	security unlock-keychain "$HOME/Library/Keychains/login.keychain-db"
fi
