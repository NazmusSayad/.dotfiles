import type { Plugin } from "@opencode-ai/plugin";

export const NoScriptExec: Plugin = async (context) => {
  return {
    "tool.execute.before": async (ctx, output) => {
      if (ctx.tool === "bash") {
        await restrictScriptExecution(output.args.command);
      }
    }
  };
};

async function restrictScriptExecution(rawCmd: unknown) {
  if (!rawCmd) return;

  const cmd = String(Array.isArray(rawCmd) ? rawCmd.join(" ") : rawCmd);
  if (!cmd.trim()) return;

  if (cmd.startsWith("node -e") || cmd.startsWith("python -c")) {
    throw new Error(
      [
        "SCRIPT OR CODE SNIPPET EXECUTION IS NOT ALLOWED for security reasons.",
        "Please USE A DIFFERENT TOOL to complete the task. If other tool can't do this, then NOTIFY THE USER.",
        "DO NOT TRY TO WRITE THIS SNIPPET AS FILE AND THEN RUN. NOTIFY THE USER how he can do it himself, but DO NOT TRY TO EXECUTE IT YOURSELF."
      ].join("\n")
    );
  }
}
