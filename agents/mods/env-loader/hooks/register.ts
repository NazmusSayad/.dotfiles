import type { Register } from 'claude-code'

const DIRENV = `cd() { builtin cd "$@" && eval "$(direnv export bash 2>/dev/null)"; }
eval "$(direnv export bash 2>/dev/null)"`

export const register: Register = on => {
  on('tool.call', { tool: 'Bash' }, ($, e, next) =>
    next({ ...e, command: `${DIRENV}\n${e.command}` }),
  )
}
