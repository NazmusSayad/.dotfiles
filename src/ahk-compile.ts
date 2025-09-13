import fs from 'fs'
import path from 'path'
import { spawnSync } from 'child_process'

const ALLOWED_SCRIPTS = ['___AHK-Macro', '___AHK-VirtualDesktop-W11']

const OUT_DIR = path.join(__dirname, '..')
const AHK_SCRIPTS = path.join(__dirname, './ahk/scripts')
const AHK2ExeBin = path.join(__dirname, './ahk/bin/Ahk2Exe.exe')
const AHKCompilerBin = path.join(__dirname, './ahk/bin/AutoHotkey64.exe')

if (!fs.existsSync(OUT_DIR)) {
  fs.mkdirSync(OUT_DIR)
}

const compiledAhkScripts = fs
  .readdirSync(AHK_SCRIPTS)
  .filter((file) => file.endsWith('.ahk'))
  .map((file) => {
    const fileName = path.basename(file, '.ahk')
    const inPath = path.join(AHK_SCRIPTS, file)
    const outPath = path.join(OUT_DIR, '___' + fileName + '.exe')
    const iconPath = path.join(AHK_SCRIPTS, fileName + '.ico')
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

function formatProgramForVBS(program: string | string[]) {
  const [exe, ...rest] = Array.isArray(program) ? program : [program]
  const literal = rest.length ? `"""${exe}"" ${rest}"` : `"""${exe}"""`

  return [
    'Set WshShell = CreateObject("WScript.Shell")',
    `WshShell.Run ${literal}, 0, False`,
    'Set WshShell = Nothing',
  ].join('\n')
}

function formatProgramForVBSAsAdmin(program: string) {
  return [
    'Set UAC = CreateObject("Shell.Application")',
    `UAC.ShellExecute "${program}", "", "", "runas", 0`,
    'Set UAC = Nothing',
  ].join('\n')
}

const vbsScript = [
  formatProgramForVBS(['cmd.exe', '/c']),

  ...[
    ...compiledAhkScripts.filter((file) =>
      ALLOWED_SCRIPTS.includes(path.basename(file, path.extname(file)))
    ),

    'C:\\Program Files\\ShareX\\ShareX.exe',
  ].map((exe) => formatProgramForVBSAsAdmin(exe)),

  formatProgramForVBS(['gpg', '--list-keys']),
]

const LAUNCH_VBS_SCRIPT = path.join(OUT_DIR, '___launch.vbs')
fs.writeFileSync(LAUNCH_VBS_SCRIPT, vbsScript.join('\n\n'))
console.log('Compilation complete')

fs.writeFileSync(
  path.join(OUT_DIR, '___task-init.xml'),
  `<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.2" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <RegistrationInfo>
    <Date>${new Date().toISOString()}</Date>
    <Author>DESKTOP-VP9L1TR\\Sayad</Author>
    <URI>\\#START</URI>
  </RegistrationInfo>
  <Triggers>
    <BootTrigger>
      <Enabled>true</Enabled>
    </BootTrigger>
    <LogonTrigger>
      <Enabled>true</Enabled>
    </LogonTrigger>
  </Triggers>
  <Principals>
    <Principal id="Author">
      <GroupId>S-1-5-32-544</GroupId>
      <RunLevel>HighestAvailable</RunLevel>
    </Principal>
  </Principals>
  <Settings>
    <MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>
    <DisallowStartIfOnBatteries>true</DisallowStartIfOnBatteries>
    <StopIfGoingOnBatteries>true</StopIfGoingOnBatteries>
    <AllowHardTerminate>true</AllowHardTerminate>
    <StartWhenAvailable>false</StartWhenAvailable>
    <RunOnlyIfNetworkAvailable>false</RunOnlyIfNetworkAvailable>
    <IdleSettings>
      <StopOnIdleEnd>true</StopOnIdleEnd>
      <RestartOnIdle>false</RestartOnIdle>
    </IdleSettings>
    <AllowStartOnDemand>true</AllowStartOnDemand>
    <Enabled>true</Enabled>
    <Hidden>true</Hidden>
    <RunOnlyIfIdle>false</RunOnlyIfIdle>
    <WakeToRun>false</WakeToRun>
    <ExecutionTimeLimit>PT72H</ExecutionTimeLimit>
    <Priority>7</Priority>
  </Settings>
  <Actions Context="Author">
    <Exec>
      <Command>${LAUNCH_VBS_SCRIPT}</Command>
    </Exec>
  </Actions>
</Task>`
)
