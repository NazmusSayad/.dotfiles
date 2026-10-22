import { spawnSync } from "node:child_process"

const IGNORED = new Set([
  "DIRENV_DIFF",
  "DIRENV_DIR",
  "DIRENV_FILE",
  "DIRENV_WATCHES",
  "XPC_SERVICE_NAME"
])

function run(command: string, cwd: string): string | undefined {
  try {
    const result = spawnSync(command, {
      env: process.env,
      cwd,
      timeout: 10000,
      shell: true
    })

    if (result.error ?? result.status !== 0) return
    const output = result.stdout.toString().trim()
    if (output === "") return
    return output
  } catch {
    return
  }
}

export default {
  id: "env-loader",
  async setup(ctx: any) {
    await ctx.shell.hook("create.before", (event: any) => {
      const targetDir: string | undefined =
        event.cwd ?? ctx.location?.directory ?? ctx.location?.project?.directory

      if (!targetDir) return

      const direnv = run("direnv export json", targetDir)
      if (direnv !== undefined) {
        let parsed: unknown
        try {
          parsed = JSON.parse(direnv)
        } catch {
          parsed = undefined
        }

        if (typeof parsed === "object" && parsed !== null) {
          for (const [key, value] of Object.entries(parsed)) {
            if (IGNORED.has(key)) continue
            if (typeof value === "string") event.env[key] = value
          }
        }
      }

      const paths = run("mise env --paths", targetDir)
      if (paths !== undefined) {
        const basePath: string | undefined = event.env.PATH ?? process.env.PATH
        const misePath = paths.split("\n").join(":")
        event.env.PATH = basePath === undefined ? misePath : `${misePath}:${basePath}`
      }

      const dotenv = run("mise env --dotenv", targetDir)
      if (dotenv !== undefined) {
        for (const line of dotenv.split("\n")) {
          const index = line.indexOf("=")
          if (index <= 0) continue
          event.env[line.slice(0, index)] = line.slice(index + 1)
        }
      }
    })
  }
}
