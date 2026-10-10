import { Plugin, type Skill } from "@opencode/plugin"
import { YAML } from "bun"
import { readFile } from "node:fs/promises"

async function loadCommands(skills: readonly Skill.Info[]) {
  const commands = []

  for (const skill of skills) {
    if (skill.path.startsWith("/builtin/")) continue

    try {
      const text = await readFile(skill.path, "utf8")
      const header = /^\uFEFF?---\r?\n([\s\S]*?)\r?\n---(?:\r?\n|$)/.exec(text)
      const frontmatter: unknown = header ? YAML.parse(header[1]) : {}

      if (typeof frontmatter !== "object" || frontmatter === null || Array.isArray(frontmatter)) {
        throw new Error("Skill frontmatter must be a YAML mapping")
      }

      let slash = true

      if ("user-invocable" in frontmatter) {
        if (typeof frontmatter["user-invocable"] !== "boolean") {
          throw new Error("user-invocable must be a boolean")
        }
        slash = frontmatter["user-invocable"]
      }

      if ("metadata" in frontmatter) {
        const metadata = frontmatter.metadata
        if (typeof metadata !== "object" || metadata === null || Array.isArray(metadata)) {
          throw new Error("Skill metadata must be a YAML mapping")
        }
        if ("opencode/slash" in metadata) {
          if (typeof metadata["opencode/slash"] !== "boolean") {
            throw new Error("opencode/slash must be a boolean")
          }
          slash = metadata["opencode/slash"]
        }
      }

      if (slash) commands.push({ name: skill.id, description: skill.description })
    } catch (error) {
      console.error(`[skill-commands] ${skill.path}:`, error)
    }
  }

  return commands
}

export default Plugin.define({
  id: "skill-commands",
  async setup(ctx) {
    const skills = await ctx.skill.list()
    let commands = await loadCommands(skills.data)

    await ctx.command.transform((editor) => {
      for (const command of commands) {
        editor.add({
          ...command,
          async execute(input) {
            await ctx.session.prompt({
              ...input.prompt,
              sessionID: input.sessionID,
              delivery: input.delivery,
              skills: [
                ...(input.prompt.skills ?? []).filter((skill) => skill.id !== command.name),
                { id: command.name }
              ]
            })
          }
        })
      }
    })

    let refreshing = false
    const timer = setInterval(async () => {
      if (refreshing) return
      refreshing = true
      try {
        const skills = await ctx.skill.list()
        const next = await loadCommands(skills.data)
        if (JSON.stringify(next) !== JSON.stringify(commands)) {
          commands = next
          await ctx.command.reload()
        }
      } catch (error) {
        console.error("[skill-commands] Refresh failed:", error)
      } finally {
        refreshing = false
      }
    }, 5000)

    return () => clearInterval(timer)
  }
})
