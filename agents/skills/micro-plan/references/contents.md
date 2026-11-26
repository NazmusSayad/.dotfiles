# Plan files

Create the plan as a folder in the current project:

```
.local/plans/<YYYY-MM-DD>-<topic>/
├── README.md
├── issue-<issue-name>.md
└── …
```

- Every issue file starts with `issue-`. `<issue-name>` is the item in a few lowercase words joined by hyphens.
- Never put numbers in the file names.

## README.md

- Title, then a line linking the sources.
- A rule for the whole plan goes here only when the user sets it for this plan (for example a naming rule). Don't carry rules over from other plans or sources.
- Then a link to each issue file: one issue per item on the user's list, in their order and wording. Nothing that isn't on the list. An issue's number (`#<N>` in replies) is its place in this list.

```markdown
# <title>

<A line linking the sources.>

<Rules the user set for this plan, if any.>

## Issues

1. [<issue title>](issue-<issue-name>.md)
2. [<issue title>](issue-<issue-name>.md)
```

## issue-<issue-name>.md

- **Title**: the issue title, in the user's wording.
- **Status**: under the title, as `**Status:** Open`. Every issue starts `Open`.
  - `Open`: not discussed yet, or still has questions.
  - `Decided`: its tasks are agreed.
  - `Skipped`: left for later, to discuss again.
  - `Dropped`: not needed, cancelled, or already fine as it is.
  - `Blocked`: waiting on something outside.
- **Details**: the facts needed to do the issue, such as what changed or was asked, the rules, and the details. Write enough that nobody has to go back to the sources.
- **Tasks**: the tasks, each an unchecked checkbox (`- [ ] …`). **Leave it empty.** The user decides the tasks; only write what they agreed to.
- **Verification**: how to check that the tasks work, each an unchecked checkbox, such as a browser check or a request to send. Leave it empty until the issue has tasks.
- **Report**: only what the user asks to know or see when the issue is done, such as a before/after video.
- **Note**: only when needed, for what the user leaves for later ("Later, not now: …", "Skipped for now, discuss later: …").
- **Context**: additional context or references, such as links to related issues, or a note that the issue is blocked by another issue. Leave it empty if there is none.

```markdown
# <issue title>

**Status:** Open

## Details

<The facts needed to do the issue.>

## Tasks

- [ ] <a task the user agreed to>

## Verification

- [ ] <how to check that a task works>

## Report

<What the user asked to know or see when the issue is done.>

## Note

<What the user left for later.>

## Context

<Additional context or references, such as related issues or the issue that blocks this one.>
```
