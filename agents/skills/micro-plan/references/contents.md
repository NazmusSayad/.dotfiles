## Plan file

Create `.local/<topic>-plan.md` in the current project.

- Title, then a line linking the sources.
- One issue per item on the user's list, in their order and wording: `## Issue N — <item>`. Nothing that isn't on the list.
- Under each issue title, its status: `**Status:** Open`. Every issue starts **Open**.
  - **Open**: not discussed yet, or still has questions.
  - **Decided**: its tasks are agreed.
  - **Skipped**: left for later, to discuss again.
  - **Dropped**: not needed, cancelled, or already fine as it is.
  - **Blocked**: waiting on something outside.
- Each issue has these sections:
  - `### Context`: the facts needed to do the issue, such as what changed or was asked, the rules, and the details. Write enough that nobody has to go back to the sources.
  - `### Tasks`: the tasks, each an unchecked checkbox (`- [ ] …`). **Leave it empty.** The user decides the tasks; only write what they agreed to.
  - `### Verification`: how to check that the tasks work, each an unchecked checkbox, such as a browser check or a request to send. Leave it empty until the issue has tasks.
  - `### Report`: only what the user asks to know or see when the issue is done, such as a before/after video.
  - `### Note`: only when needed, for what the user leaves for later ("Later, not now: …", "Skipped for now, discuss later: …").
  - `### Context`: Additional context, such as links to related issues, or a note that the issue is blocked by another issue. Leave it empty if there is no additional context.
- A rule for the whole plan goes at the top only when the user sets it for this plan (for example a naming rule). Don't carry rules over from other plans or sources. 