# Plan files

Create the plan as a folder in the current project:

```
.local/plans/<YYYY-MM-DD>-<topic>/
├── plan.md
├── <issue-name>.md
└── …
```

- `<issue-name>` is the item in a few lowercase words joined by hyphens.
- Never put numbers in the file names.

## plan.md

- Title, then a line linking the sources.
- A rule for the whole plan goes here only when the user sets it for this plan (for example a naming rule). Don't carry rules over from other plans or sources.
- Then a link to each issue file: one issue per item on the user's list, in their order and wording. Nothing that isn't on the list. An issue's number (`#<N>` in replies) is its place in this list.

## <issue-name>.md

- Title: the item, in the user's wording.
- Under the title, its status: `**Status:** Open`. Every issue starts **Open**.
  - **Open**: not discussed yet, or still has questions.
  - **Decided**: its tasks are agreed.
  - **Skipped**: left for later, to discuss again.
  - **Dropped**: not needed, cancelled, or already fine as it is.
  - **Blocked**: waiting on something outside.
- These sections:
  - `## Details`: the facts needed to do the issue, such as what changed or was asked, the rules, and the details. Write enough that nobody has to go back to the sources.
  - `## Tasks`: the tasks, each an unchecked checkbox (`- [ ] …`). **Leave it empty.** The user decides the tasks; only write what they agreed to.
  - `## Verification`: how to check that the tasks work, each an unchecked checkbox, such as a browser check or a request to send. Leave it empty until the issue has tasks.
  - `## Report`: only what the user asks to know or see when the issue is done, such as a before/after video.
  - `## Note`: only when needed, for what the user leaves for later ("Later, not now: …", "Skipped for now, discuss later: …").
  - `## Context`: Additional context, such as links to related issues, or a note that the issue is blocked by another issue. Leave it empty if there is no additional context.
