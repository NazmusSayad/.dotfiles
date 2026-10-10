import type { EngineInterface, Register } from 'claude-code'

let branch = ''

async function refresh($: EngineInterface) {
  const { exitCode, stdout } = await $.process.run(['git', 'branch', '--show-current'])
  branch = exitCode === 0 ? stdout.trim() : ''
  $.ui.invalidate('ui.render')
}

export const register: Register = on => {
  on('session.start', async ($, e, next) => {
    await refresh($)

    return next(e)
  })

  on('turn.complete', async ($, e, next) => {
    await refresh($)

    return next(e)
  })

  on('ui.render', { component: 'PromptHint' }, ($, e, next) => {
    if (branch === '') {
      return next(e)
    }

    return next({ ...e, props: { ...e.props, tail: ` · ${branch}` } })
  })
}
