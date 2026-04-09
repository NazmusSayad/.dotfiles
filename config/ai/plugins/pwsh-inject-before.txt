import type { Plugin } from "@opencode-ai/plugin"

export const PwshInjectPlugin: Plugin = async () => {
  return {
    "tool.execute.before": async (input, output) => {
      const isRunningPwsh =
        input.tool === "bash" && typeof output?.args?.command === "string"

      if (isRunningPwsh) {
        output.args.command = [
          ". $PROFILE; # injected by plugin; ignore",
          output.args.command
        ].join("\n\n")
      }
    }
  }
}
