import { spawnSync } from "node:child_process"

const IGNORED = new Set([
  "DIRENV_DIFF",
  "DIRENV_DIR",
  "DIRENV_FILE",
  "DIRENV_WATCHES",
  "XPC_SERVICE_NAME"
])

export default {
  id: "env-loader",
  async setup(ctx: any) {
    await ctx.shell.hook("create.before", (event: any) => {
      let raw: string
      try {
        const result = spawnSync("direnv", ["export", "json"], {
          env: process.env,
          cwd: event.cwd,
          timeout: 10000,
          shell: true
        })

        if (result.error ?? result.status !== 0) return
        raw = result.stdout.toString().trim()
        if (raw === "") return
      } catch {
        return
      }

      let parsed: unknown
      try {
        parsed = JSON.parse(raw)
      } catch {
        return
      }

      if (typeof parsed !== "object" || parsed === null) return

      for (const [key, value] of Object.entries(parsed)) {
        if (IGNORED.has(key)) continue
        if (typeof value === "string") event.env[key] = value
      }
    })
  }
}
