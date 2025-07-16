import fs from 'fs'
import path from 'path'
import { execSync } from 'child_process'

function findBinPath(bin: string): string | null {
  try {
    const stdout = execSync(
      `${process.platform === 'win32' ? 'where ' : 'which'} ${bin}`,

      {
        stdio: ['ignore', 'pipe', 'ignore'],
      }
    ).toString()
    return (
      stdout
        .split('\n')
        .map((line) => line.trim())
        .filter(Boolean)[0] || null
    )
  } catch {
    return null
  }
}

const EDITORS = [
  {
    name: 'VSCode',
    alias: 'code',
    extensionsPath: '../../resources/app/extensions',
  },
  {
    name: 'VSCode Insiders',
    alias: 'code-insiders',
    extensionsPath: '../../resources/app/extensions',
  },
  {
    name: 'Cursor',
    alias: 'cursor',
    extensionsPath: '../../../../resources/app/extensions',
  },
]

const matchedEditors = EDITORS.map((a) => ({
  ...a,
  path: findBinPath(a.alias),
})).filter((a) => a.path)

function getFiles(root: string): string[] {
  const stat = fs.readdirSync(root, { withFileTypes: true })

  const files = stat.map((dirent) => {
    const abc = path.join(root, dirent.name)
    return dirent.isFile() ? abc : getFiles(abc)
  })

  return files.flat()
}

matchedEditors.forEach((editor) => {
  const getExtensionsPath = path.join(editor.path, editor.extensionsPath)

  if (!fs.existsSync(getExtensionsPath)) {
    return console.warn(`Extensions path not found for ${editor.name}`, {
      path: editor.path,
      extPath: editor.extensionsPath,
      resolvedPath: getExtensionsPath,
    })
  }

  const foundFiles = getFiles(getExtensionsPath)
    .filter((file) => file.toLowerCase().endsWith('.code-snippets'))
    .sort((a, b) => path.basename(a).length - path.basename(b).length)

  foundFiles.forEach((file) => {
    fs.writeFileSync(file, '{}')
    console.log(`${editor.name} Cleared: ${path.basename(file)}`)
  })
})

console.log('Done!')
