import type { EngineInterface, Register } from 'claude-code'

let git = ''

async function refreshGit($: EngineInterface) {
  const { exitCode, stdout } = await $.process.run(['git', 'status', '--porcelain=v2', '--branch'])
  git = ''
  if (exitCode === 0) {
    const lines = stdout.split('\n')
    const head = lines.find(line => line.startsWith('# branch.head '))?.slice(14) ?? ''
    const [ahead, behind] = (lines.find(line => line.startsWith('# branch.ab '))?.slice(12) ?? '+0 -0').split(' ')
    const isDirty = lines.some(line => line !== '' && !line.startsWith('#'))
    git = `${head}${isDirty ? '*' : ''}${ahead !== '+0' ? '↑' : ''}${behind !== '-0' ? '↓' : ''}`
  }

  $.ui.invalidate('ui.render')
}

export const register: Register = on => {
  on('session.start', async ($, e, next) => {
    await refreshGit($)
    $.clock.every(5000, () => refreshGit($))

    return next(e)
  })

  on('turn.complete', async ($, e, next) => {
    await refreshGit($)

    return next(e)
  })

  on('ui.render', { component: 'PromptHint' }, ($, e, next) => {
    if (git === '') {
      return next(e)
    }

    return next({ ...e, props: { ...e.props, tail: git } })
  })
}
