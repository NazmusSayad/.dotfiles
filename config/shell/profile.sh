[[ -f ~/.path ]] && export PATH="$PATH:$(paste -sd ':' ~/.path)"

eval "$(/opt/homebrew/bin/brew shellenv bash)"
eval "$(mise env --shell bash)"

[[ -f ~/.env ]] && eval "$(dotsh bash "$(cat ~/.env)")"

export RUST_BACKTRACE=1
export NODE_NO_WARNINGS=1

export OPENCODE_SHELL="/opt/homebrew/bin/bash"