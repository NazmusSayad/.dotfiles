The PR body is a briefing, not a lab notebook. A reviewer who has the diff open should learn why the change exists, what it changes, what it leaves out, what could break, and how you know it works, in under a minute.

Write every piece of text, including headings and labels, in strict ASD-STE100 Simplified Technical English, with no jargon. Use the project's own names for things, and name a real symbol or path only when it carries the change. Give the consequence, not the mechanics. Skip preambles. Keep the body under about 40 lines and move longer background, full output, or logs into a collapsed `<details>` section at the end.

Base every claim on code, output, or results you actually inspected. Never invent intent, impact, or test results. When you can't tell, say so.

If the repository has a PR template or title convention, follow it and fit these sections into it.

## Title

A short instruction that says what the change does, with no trailing period.

## Body

Use these sections in this order, each under a `##` heading. Drop Other options when it has nothing to say.

- **Why:** the problem and the approach in one to three short sentences. Link the issue when there is one.
- **Before and after:** the change shown as a before and after. Each side can be anything that makes the difference clear: a single line of text, output, a screenshot, or a video recording.
- **Risk:** one or two sentences on who or what the change touches, what could break, and why it is safe or risky. Only say something breaks when you can name what breaks. Say whether reverting the PR fully undoes it, and call out anything a revert can't undo.
- **Evidence:** proof that the change works, not a claim that tests pass. Show the exact test that failed before and passes now, as pseudocode, or the output that changed. For a performance change, give one number with its unit as `before → after`. Say what you couldn't check.
- **What changed:** what the PR changes and what it leaves out on purpose, in one to three short bullets, not a file list. Name both sides of a rename or move, and call out changes the task didn't ask for.
- **Other options:** only the alternatives a reviewer would otherwise ask about.

## Attachments

Add screenshots whenever the change is visible and the environment can take them. They are the strongest evidence. Crop to the part that changed and place them next to the text they support.

Add a video recording when screenshots can't show the change, because it moves or spans several steps. Start from a screen the reviewer knows and show only the steps that matter.

Never add attachment files to git. `gh` supports attachments (gh 2.99.0 or later): reference each file by its path in the body and pass the same path to `--attach` on `gh pr create` or `gh pr edit`. `gh` uploads the file and replaces the path with its link, so it shows where you placed it.

## Diagrams

Add a Mermaid diagram when the change affects how parts connect or the order in which things happen. GitHub draws it in the PR. Keep it small, show only what changed, and show before and after when the flow itself changed. Place it next to the text it explains.

## Updating

When the PR already exists, rewrite the title and body to describe the whole PR as it is now, without update notes. Keep attachments, diagrams, and evidence that still hold, and replace the ones that don't.
