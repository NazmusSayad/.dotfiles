eval "$(/opt/homebrew/bin/brew shellenv bash)"
eval "$(mise env --shell bash)"

[[ -f ~/.path ]] && export PATH="$PATH:$(paste -sd ':' ~/.path)"
[[ -f ~/.env ]] && eval "$(dotsh bash "$(cat ~/.env)")"

export OPENCODE_DISABLE_CLAUDE_CODE=1