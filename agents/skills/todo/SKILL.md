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
- One issue per item on the user's list, in their order and wording: `## Issue N — <item>`.
- Related items not on the list go last, as extra issues marked "(not in the list)".
- Each issue has these sections:
  - `### Context`: the facts needed to do the issue, such as what changed or was asked, the rules, and the details. Write enough that nobody has to go back to the sources.
  - `### Tasks`: the tasks, each an unchecked checkbox (`- [ ] …`). **Leave it empty.** The user decides the tasks; only write what they agreed to.
  - `### Verification`: how to check that the tasks work, each an unchecked checkbox, such as a browser check or a request to send. Leave it empty until the issue has tasks.
  - `### Note`: only when needed, for items that are skipped or postponed ("Later, not now: …", "Skipped for now, discuss later: …").
- If the user sets a rule for the whole plan (for example a naming rule), put it at the top.

## First reply

Keep it short:
- where the plan file is;
- one or two lines per issue on where things stand;
- a ⚠️ for anything that will break;
- the open questions.

## Discussing the issues

Go one issue per turn, in order. Keep each turn short. Long, technical replies frustrate the user.

```
**Issue N: <title>**

<What the item is about and where things stand now. As short as possible without losing any important detail.>

📋 **TODO:**
- **<where>:** <what changes, as the user will see it>

💡 **Solutions:**
- **<topic>:** <the solution, and why it fits>

❓ **Questions:**
- **<topic>:** <the question, and your recommendation>

🚫 **Blockers:**
- **<topic>:** <what stops the work, and who or what can unblock it>

⚠️ **Risks:**
- **<topic>:** <a major problem the change could cause, and how to avoid it>
```

- **Solutions** are how to meet the item. If something that already exists meets it, that is the solution; point to it instead of inventing a new one. If the solution is obvious, give the obvious one.
- **Questions** are real decisions the user has to make between two or more options. Never ask the user to approve the TODO or the solution, such as "Go ahead with this?" or "OK with this?", in **Questions** or anywhere else. The user will say if they disagree.
- **Blockers** are things that stop the work, such as a missing API, access or information.
- **Risks** are major problems the change could cause, such as breaking an existing feature or losing data.
- Show **Solutions**, **Questions**, **Blockers** and **Risks** only when they have content.
- Leave out technical details unless the user needs them to decide.
- If something already works, say so and suggest nothing, or only a quick check.
- Answer the user's questions directly, after checking the code.

After each answer:
1. Write exactly what the user decided into that issue's **Tasks** as unchecked checkboxes, and nothing more. Postponed or skipped items go under **Note**.
2. Write how to check those tasks under **Verification**.
3. Say in one line what you added. If this issue still has open questions or blockers, present this issue again in the format above. Present the next issue only when nothing in this issue is open, or when the user says to move on.

When every issue is done, list the notes that were left for later.

Don't implement anything until the user says to start.
