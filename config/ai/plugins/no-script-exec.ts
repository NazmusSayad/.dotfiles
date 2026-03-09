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
        "Script or code snippet execution is not allowed for security reasons.",
        "Please USE A DIFFERENT TOOL to complete the task. If there isn't any other way, then notify the user.",
        "DO NOT WRITE A SCRIPT OR CODE SNIPPET TO EXECUTE THE TASK. INSTEAD, EXPLAIN TO THE USER HOW THEY CAN DO IT THEMSELVES.",
        "DO NOT TRY TO EXECUTE SCRIPTS OR CODE SNIPPETS DIRECTLY."
      ].join("\n")
    );
  }
}
