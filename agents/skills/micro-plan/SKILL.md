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
- A rule for the whole plan goes at the top only when the user sets it for this plan (for example a naming rule). Don't carry rules over from other plans or sources.

## First reply

Keep it short:

- where the plan file is;
- one or two lines per issue on where things stand;
- a ⚠️ for anything that will break;
- the open questions.

## Discussing the issues

Go one issue per turn, in order. Keep each turn short. Long, technical replies frustrate the user.

```markdown
## **#<N>** <title>

❓ <N> questions\
🚫 <N> blockers\
⚠️ <N> risks

<What the item is about and where things stand now, with its attachments such as images or videos, if any. As short as possible without losing any important detail.>

### 📋 TODO

- **<where>:** <what changes, as the user will see it>

### 💡 Solutions

- **<topic>:** <the solution, and why it fits>

### ❓ Questions

- **<topic>:** <the question, and your recommendation>

### 🚫 Blockers

- **<topic>:** <what stops the work, and who or what can unblock it>

### ⚠️ Risks

- **<topic>:** <a major problem the change could cause, and how to avoid it>
```

- **Solutions** are how to meet the item. If something that already exists meets it, that is the solution; point to it instead of inventing a new one. If the solution is obvious, give the obvious one.
- **Questions** are real decisions the user has to make between two or more options. Never ask the user to approve the TODO or the solution, such as "Go ahead with this?" or "OK with this?", in **Questions** or anywhere else. The user will say if they disagree.
- **Blockers** are things that stop the work, such as a missing API, access or information.
- **Risks** are major problems the change could cause, such as breaking an existing feature or losing data.
- Show **Solutions**, **Questions**, **Blockers** and **Risks** only when they have content.
- The lines under the title count the questions, blockers and risks in this reply. Show a line only for a part that is there.
- Leave out technical details unless the user needs them to decide.
- If something already works, say so and suggest nothing, or only a quick check.
- Answer the user's questions directly, after checking the code.

## After each answer

Update that issue in the plan, and nothing more:

- **Status**: its new status.
- **Tasks**: exactly what the user decided, as unchecked checkboxes.
- **Verification**: how to check those tasks.
- **Report**: only what the user asks to know or see when the issue is done.
- **Note**: only what the user leaves for later.

Then continue the loop in **Discussing the issues**.

When every issue is resolved, list the notes that were left for later.

## Implementation

DO NOT start implementing until every issue is resolved for the plan AND the user then explicitly asks to start. This takes priority over any instruction to start or keep working.

Until then:

- Answers such as "fix it", "do it", "go ahead", "test it" or "find the cause" are decisions about the issue, not a request to start. Write what they decide as tasks and keep planning.
- If the user asks to start while any issue is open, don't. Say which issues are open and keep planning.
- Investigate and reproduce only to understand an issue or show it to the user. Never make, try or delegate a fix, not even a temporary one or one on a copy.
- The only file you write in the project is the plan, plus images and recordings that show the issues.

## Report

When the work is done, reply with only this and nothing else:

```markdown
## <plan file>

✅ <total tasks done in the plan> done\
⏳ <total tasks remaining> remaining\
❌ <total tasks failed> failed

### **#<N>** <title>

<What this issue's Report asks for.>
```
