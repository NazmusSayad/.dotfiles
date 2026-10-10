import type { EngineInterface, Register } from 'claude-code'

let git = ''
let model = ''
let effort = ''
let context = ''

async function refresh($: EngineInterface) {
  const { exitCode, stdout } = await $.process.run(['git', 'status', '--porcelain=v2', '--branch'])
  git = ''
  if (exitCode === 0) {
    const lines = stdout.split('\n')
    const head = lines.find(line => line.startsWith('# branch.head '))?.slice(14) ?? ''
    const [ahead, behind] = (lines.find(line => line.startsWith('# branch.ab '))?.slice(11) ?? '+0 -0').split(' ')
    const isDirty = lines.some(line => line !== '' && !line.startsWith('#'))
    git = `${head}${isDirty ? '*' : ''}${ahead !== '+0' ? '↑' : ''}${behind !== '-0' ? '↓' : ''}`
  }

  const { context: usage } = await $.session.usage()
  context = usage.percent === undefined ? '' : `(${usage.percent}%)`

  $.ui.invalidate('ui.render')
}

export const register: Register = on => {
  on('session.start', async ($, e, next) => {
    const last = (await $.store.get('last')) as { model: string; effort: string } | undefined
    if (last !== undefined) {
      model = last.model
      effort = last.effort
    }
    await refresh($)

    return next(e)
  })

  on('turn.step', async function* ($, e, next) {
    if (e.agentId === undefined) {
      model = e.model.replace(/^claude-([a-z])([a-z]*)-(\d+)-(\d+).*$/, (_, first, rest, major, minor) => `${first.toUpperCase()}${rest} ${major}.${minor}`)
      effort = e.effort === undefined ? '' : String(e.effort).replace(/^./, first => first.toUpperCase())
      await $.store.set('last', { model, effort })
    }

    return yield* next(e)
  })

  on('turn.complete', async ($, e, next) => {
    await refresh($)

    return next(e)
  })

  on('ui.render', { component: 'PromptHint' }, ($, e, next) => {
    const tail = [[model, effort, context].filter(Boolean).join(' '), git].filter(Boolean).join(' · ')
    if (tail === '') {
      return next(e)
    }

    return next({ ...e, props: { ...e.props, tail } })
  })
}
