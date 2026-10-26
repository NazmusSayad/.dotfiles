---
name: todo
description: Plans a list of items step by step with the user.

disable-model-invocation: true
metadata: { opencode/autoinvoke: false }
---

## Input

The user gives a short list of items. They may point to a repo, docs, files or a past session for context.

## Research

1. Gather the sources the user points to. Pull the repos they name. Read their docs first, and read code only where the docs are missing or unclear. If the user uploads files, compare them with the repo copies and use the newest.
2. For each item, find the code or area it touches. Note what already works, what is missing, and what will break. Check the code before saying something already works.
3. Also note related things that aren't on the user's list.

## Plan file

Create `docs/<topic>-plan.md` in the current project.

- Title, then a line linking the sources.
- One step per item on the user's list, in their order and wording: `## Step N — <item>`.
- Related items not on the list go last, as extra steps marked "(not in the list)".
- Each step has these sections:
  - `### Context`: the facts needed to do the step, such as what changed or was asked, the rules, and the details. Write enough that nobody has to go back to the sources.
  - `### Tasks`: the tasks, each an unchecked checkbox (`- [ ] …`). **Leave it empty.** The user decides the tasks; only write what they agreed to.
  - `### Note`: only when needed, for items that are skipped or postponed ("Later, not now: …", "Skipped for now, discuss later: …").
- If the user sets a rule for the whole plan (for example a naming rule), put it at the top.

## First reply

Keep it short:
- where the plan file is;
- one or two lines per step on where things stand;
- a ⚠️ for anything that will break;
- the open questions.

## Discussing the steps

Go one step per turn, in order. Keep each turn short. Long, technical replies frustrate the user.

```
**Step N: <title>**

<What the item is about and where things stand now. As short as possible without losing any important detail.>

**TODO:**
- **<where>:** <what changes, as the user will see it>

**Questions:**
- **<topic>:** <the question, the options, and which one you'd pick>

**Blockers:**
- **<topic>:** <what stops the work, and who or what can unblock it>
```

- **Questions** are decisions the user has to make. **Blockers** are things that stop the work, such as a missing API, access or information.
- Show **Questions** and **Blockers** only when they have content. Don't end with "OK with this?".
- Leave out technical details unless the user needs them to decide.
- If something already works, say so and suggest nothing, or only a quick check.
- Answer the user's questions directly, after checking the code.

After each answer:
1. Write exactly what the user decided into that step's **Tasks** as unchecked checkboxes, and nothing more. Postponed or skipped items go under **Note**.
2. Say in one line what you added, then present the next step.

When every step is done, list the notes that were left for later.

Don't implement anything until the user says to start.
