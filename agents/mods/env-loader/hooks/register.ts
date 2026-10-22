import type { Register } from 'claude-code'

const ENV = `load_env() {
  eval "$(direnv export bash 2>/dev/null)"
  local paths="$(mise env --paths 2>/dev/null | paste -sd: -)"
  [ -n "$paths" ] && export PATH="$paths:$PATH"
  while IFS= read -r line; do export "$line"; done < <(mise env --dotenv 2>/dev/null)
}
cd() { builtin cd "$@" && load_env; }
load_env`

export const register: Register = on => {
  on('tool.call', { tool: 'Bash' }, ($, e, next) =>
    next({ ...e, command: `${ENV}\n${e.command}` }),
  )
}
