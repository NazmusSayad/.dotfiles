import { YAML } from "bun"
import { readFile } from "node:fs/promises"

type Prompt = {
  text: string
  skills?: { id: string }[]
}

type Invocation = {
  sessionID: string
  prompt: Prompt
  delivery: "steer" | "queue"
}

type Context = {
  skill: {
    list(): Promise<{
      data: { id: string; description?: string; path: string }[]
    }>
  }
  command: {
    transform(callback: (editor: {
      add(command: {
        name: string
        description?: string
        execute(input: Invocation): Promise<void>
      }): void
    }) => void): Promise<unknown>
    reload(): Promise<void>
  }
  session: {
    prompt(input: Prompt & {
      sessionID: string
      delivery: "steer" | "queue"
    }): Promise<unknown>
  }
}

async function loadCommands(ctx: Context) {
  const skills = await ctx.skill.list()
  const commands = []

  for (const skill of skills.data) {
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

export default {
  id: "skill-commands",
  async setup(ctx: Context) {
    let commands = await loadCommands(ctx)

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
        const next = await loadCommands(ctx)
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
}
