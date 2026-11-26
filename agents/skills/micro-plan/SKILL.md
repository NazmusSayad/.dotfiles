---
name: micro-plan
description: Plans a list of items step by step with the user.

disable-model-invocation: true
metadata: { opencode/autoinvoke: false }
---

This skill turns a list of items into a plan by discussing them with the user one issue at a time and writing down what they decide.

## Input

The user gives a short list of items. They may point to a repo, docs, files or a past session for context.

## Research

1. Gather the sources the user points to. Pull the repos they name. Read their docs first, and read code only where the docs are missing or unclear. If the user uploads files, compare them with the repo copies and use the newest.
2. For each item, find the code or area it touches. Note what already works, what is missing, and what will break. Check the code before saying something already works.

## Plan files

Before you write any plan file, read `references/contents.md`. Create and update the plan files exactly as it describes: the same folder, file names, sections and order.

## Communication

Before your first reply, read `references/communication.md`. Write every reply to the user exactly as it describes: the first reply, each issue, and the final report.

## After each answer

Update that issue's file, and nothing more:

- **Status**: its new status.
- **Tasks**: exactly what the user decided, as unchecked checkboxes.
- **Verification**: how to check those tasks.
- **Report**: only what the user asks to know or see when the issue is done.
- **Note**: only what the user leaves for later or mentions in passing.

Then continue the loop in **Discussing the issues** in `references/communication.md`.

When every issue is resolved, list the notes that were left for later.

## Implementation

DO NOT start implementing until no issue is `Open` AND the user then explicitly asks to start. This takes priority over any instruction to start or keep working.

Until then:

- Answers such as "fix it", "do it", "go ahead", "test it" or "find the cause" are decisions about the issue, not a request to start. Write what they decide as tasks and keep planning.
- If the user asks to start while any issue is `Open`, don't. Say which issues are open and keep planning.
- Investigate and reproduce only to understand an issue or show it to the user. Never make, try or delegate a fix, not even a temporary one or one on a copy.
- The only files you write in the project are the plan files, plus images and recordings that show the issues.

Once the user asks to start, implement the `Decided` issues one after another, or in parallel with subagents when you can. As you go:

- Capture what each issue's Report needs.
- Update each issue's file: tick each task that is done (`- [x] …`), and mark a task that fails with ❌ and why (`- [ ] ❌ … — <why>`).
- After an issue's tasks, run its Verification. Tick each check that passes, and mark one that fails the same way as a failed task.
