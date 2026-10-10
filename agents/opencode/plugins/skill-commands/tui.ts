import { Plugin } from "@opencode/plugin/tui"

export default Plugin.define({
  id: "skill-commands.tui",
  setup(ctx) {
    return ctx.data.on("command.updated", async (event) => {
      if (event.location) return
      const location = ctx.location ?? ctx.data.location.default()
      ctx.data.location.command.invalidate(location)
      try {
        await ctx.data.location.command.sync(location)
      } catch (error) {
        console.error("[skill-commands] Command suggestions refresh failed:", error)
      }
    })
  }
})
