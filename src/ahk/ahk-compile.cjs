const fs = require('fs')
const path = require('path')
const { spawnSync } = require('child_process')

const AHKOutDir = path.join(__dirname, '../bin')
const AHKIons = path.join(__dirname, './icons')
const AHKScripts = path.join(__dirname, './scripts')
const AHK2ExeBin = path.join(__dirname, './bin/Ahk2Exe.exe')
const AHKCompilerBin = path.join(__dirname, './bin/AutoHotkey64.exe')

if (fs.existsSync(AHKOutDir)) {
  fs.rmSync(AHKOutDir, { recursive: true, force: true })
}
fs.mkdirSync(AHKOutDir)

const compiledAhkScripts = fs
  .readdirSync(AHKScripts)
  .filter((file) => file.endsWith('.ahk'))
  .map((file) => {
    const fileName = path.basename(file, '.ahk')
    const inPath = path.join(AHKScripts, file)
    const outPath = path.join(AHKOutDir, fileName + '.exe')
    const iconPath = path.join(AHKIons, fileName + '.ico')
    const icoExists = fs.existsSync(iconPath)

    const spawnArgs = [
      '/base',
      AHKCompilerBin,
      '/in',
      inPath,
      '/out',
      outPath,
      icoExists && ['/icon', iconPath],
    ]

    spawnSync(AHK2ExeBin, spawnArgs.flat().filter(Boolean), {
      stdio: 'inherit',
      shell: true,
    })

    console.log(`Compiled: ${file}`)
    return outPath
  })

const vbsScript = [
  ['cmd.exe', '/c'],
  ...compiledAhkScripts,
  ['C:\\Program Files\\ShareX\\ShareX.exe'],
  ['gpg', '--list-keys'],
].map((program) =>
  [
    'Set WshShell = CreateObject("WScript.Shell")',
    formatProgramForVBS(program),
    'Set WshShell = Nothing',
  ].join('\n')
)

function formatProgramForVBS(program) {
  const [exe, ...rest] = Array.isArray(program) ? program : [program]
  const literal = rest.length ? `"""${exe}"" ${rest}"` : `"""${exe}"""`
  return `WshShell.Run ${literal}, 0, False`
}

fs.writeFileSync(path.join(AHKOutDir, 'Launch.vbs'), vbsScript.join('\n\n'))
console.log('Compilation complete')
