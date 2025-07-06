import fs from 'fs'
import path from 'path'
import { addToPath, writeEnv } from './env.mjs'

const MSYS_PATH = path.resolve('C:\\msys64')
const NSSWITCH_CONFIG_PATH = path.resolve(MSYS_PATH, 'etc/nsswitch.conf')
const MSYS_INIS_PATH = [
  'msys2.ini',
  'clang32.ini',
  'clang64.ini',
  'clangarm64.ini',
  'mingw32.ini',
  'mingw64.ini',
  'ucrt64.ini',
]
const MSYS_BINS = [
  'usr/bin',
  'mingw64/bin',
  'mingw32/bin',
  'ucrt64/bin',
  'clang32/bin',
  'clang64/bin',
  'clangarm64/bin',
]

for (const ini of MSYS_INIS_PATH) {
  const iniPath = path.resolve(MSYS_PATH, ini)
  const isExist = fs.existsSync(iniPath)
  if (!isExist) {
    console.log(`File not found: ${iniPath}`)
    continue
  }

  const iniContent = fs.readFileSync(iniPath, 'utf-8')
  fs.writeFileSync(
    iniPath,
    iniContent.replace(/^\#(?=MSYS2_PATH_TYPE\=inherit)/gm, '')
  )
  console.log(`Updated: ${ini}`)
}

const nsswitchConfig = fs.readFileSync(NSSWITCH_CONFIG_PATH, 'utf-8')
fs.writeFileSync(
  NSSWITCH_CONFIG_PATH,
  nsswitchConfig.replaceAll(
    /(?<=(?:db_home|db_shell|db_gecos): ).+/g,
    'windows'
  )
)
console.log('Updated: nsswitch.conf')

writeEnv('MSYS2_PATH_TYPE', 'inherit', 'Machine')
addToPath(
  'Machine',
  ...MSYS_BINS.map((bin) => path.join(MSYS_PATH, bin)).filter(
    (bin) => fs.existsSync(bin) && fs.statSync(bin).isDirectory()
  )
)
