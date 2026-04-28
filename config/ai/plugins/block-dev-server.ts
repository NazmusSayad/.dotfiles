import { Plugin } from "@opencode-ai/plugin"

function npmRunCommand(commands: string): string[] {
  return [
    `bun ${commands}`,
    `pnpm ${commands}`,
    `yarn ${commands}`,
    `npm run ${commands}`,
    `bun run ${commands}`,
    `yarn run ${commands}`,
    `pnpm run ${commands}`
  ]
}

function npxCommand(commands: string): string[] {
  return [
    commands,
    `npx ${commands}`,
    `bunx ${commands}`,
    `pnpx ${commands}`,
    `yarnx ${commands}`,
    `npm exec ${commands}`,
    `bun exec ${commands}`,
    `pnpm exec ${commands}`,
    `yarn exec ${commands}`,
    `pnpm dlx ${commands}`,
    `yarn dlx ${commands}`
  ]
}

const BLOCKED_PATTERNS: string[] = [
  ...npmRunCommand("dev"),
  ...npmRunCommand("start"),
  ...npmRunCommand("serve"),

  ...npxCommand("vite"),
  ...npxCommand("next"),
  ...npxCommand("nodemon"),
  ...npxCommand("tsx watch"),
  ...npxCommand("tsx --watch"),
  ...npxCommand("ts-node-dev"),
  ...npxCommand("react-scripts start"),

  ...npxCommand("serve"),
  ...npxCommand("http-server"),
  ...npxCommand("live-server"),

  ...npxCommand("remix dev"),
  ...npxCommand("astro dev"),
  ...npxCommand("expo start"),
  ...npxCommand("svelte-kit dev"),

  ...npxCommand("hexo server"),
  ...npxCommand("gatsby develop"),
  ...npxCommand("solid-start dev"),
  ...npxCommand("quartz build --serve"),
  ...npxCommand("vue-cli-service serve"),

  "node server.js",
  "node --watch",
  "pm2 start",

  "parcel",
  "webpack serve",
  "rollup -c -w",
  "esbuild --watch"
]

export const BlockDevServer: Plugin = async () => {
  return {
    "tool.execute.before": async (input: any, output: any) => {
      if (input.tool === "bash") {
        const command: string = (output.args.command ?? "").trim()

        for (const p of BLOCKED_PATTERNS) {
          const cmds = command
            .split(/\s+/gim)
            .join(" ")
            .split(/;|&|&&/gim)

          if (cmds.some((c) => c.startsWith(p))) {
            throw new Error(
              [
                `Command usage restricted: "${command}".`,
                `You should NOT run this command. DO NOT TRY TO BYPASS THIS RESTRICTION!`,
                `Instead ask the user to run "${command}" or continue other tasks.`
              ].join("\n")
            )
          }
        }
      }
    }
  }
}
