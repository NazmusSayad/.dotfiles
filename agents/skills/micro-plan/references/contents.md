Create the plan as a folder in the current project:

```
.local/plans/<YYYY-MM-DD>-<topic>/
├── README.md
└── issues/
    ├── <issue-name>/
    │   ├── README.md
    │   └── <images and recordings for this issue>
    └── …
```

- `<issue-name>` is the issue title in a few lowercase words joined by hyphens.
- Never put numbers in file or folder names.
- Put each issue's images and recordings in its folder.

## README.md

The overview of the plan. Keep it up to date with the user's feedback.

- Title, then a line linking the sources.
- A rule for the whole plan goes here only when the user sets it for this plan (for example a naming rule). Don't carry rules over from other plans or sources.
- Then a link to each issue's `README.md`: one issue per item on the user's list, in their order and wording. Nothing that isn't on the list. An issue's number (`#<N>` in replies) is its place in this list.

```markdown
# <title>

<A line linking the sources.>

<Rules the user set for this plan, if any.>

## Issues

1. [<issue title>](issues/<issue-name>/README.md)
2. [<issue title>](issues/<issue-name>/README.md)
```

## issues/<issue-name>/README.md

The plan for one issue. Keep it up to date with the user's feedback and the progress of its tasks.

- **Title**: the issue title, in the user's wording.
- **Status**: under the title, as `**Status:** Open`. Every issue starts `Open`.
  - `Open`: not discussed yet, or still has questions.
  - `Decided`: its questions and blockers are resolved.
  - `Skipped`: left for later, to discuss again.
  - `Dropped`: not needed, cancelled, or already fine as it is.
  - `Blocked`: waiting on something outside.
- **Details**: the facts needed to do the issue, such as what changed or was asked, the rules, and the details. Write enough that nobody has to go back to the sources.
- **Tasks**: the tasks, each an unchecked checkbox (`- [ ] …`). **Leave it empty until the issue is `Decided`.** Then write its TODO as the user last saw it, with what they decided, and nothing else.
- **Verification**: how to check by hand that the tasks really work, each an unchecked checkbox, such as using the change in the browser or sending a real request in the proper environment. It is not a static check such as lint, formatting or type checks. Write it once the issue has tasks.
- **Report**: only what the user asks to know or see when the issue is done, such as a before/after video.
- **Note**: only when needed, for what the user leaves for later or mentions in passing ("Later, not now: …", "Skipped for now, discuss later: …", "Waits on #2").
- **Context**: additional context or references, such as links to related issues, or a note that the issue is blocked by another issue. Leave it empty if there is none.

```markdown
# <issue title>

**Status:** Open

## Details

<The facts needed to do the issue.>

## Tasks

- [ ] <a task the user agreed to>

## Verification

- [ ] <a hands-on check that a task works, such as in the browser>

## Report

<What the user asked to know or see when the issue is done.>

## Note

<What the user left for later or mentioned in passing.>

## Context

<Additional context or references, such as related issues or the issue that blocks this one.>
```
