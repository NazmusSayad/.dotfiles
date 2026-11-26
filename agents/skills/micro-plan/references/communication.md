## First reply

Keep it short:

- where the plan is;
- one or two lines per issue on where things stand;
- a ⚠️ for anything that will break;
- the open questions.

## Discussing the issues

Go one issue per turn, in order. Keep each turn short. Long, technical replies frustrate the user.

The first time you present an issue, or when the user asks to see it again, use this structure:

```markdown
<One line on what you did with the user's last reply, if you have anything to say about it.>

---

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

While you and the user are still discussing that issue, use this one instead, with only what is still open:

```markdown
<One line on what you did with the user's last reply, if you have anything to say about it.>

---

## 💬 #<N> <title>

### ❓ Questions (<open questions> of <all questions this issue had> remaining)

- **<topic>:** <the question, and your recommendation>

### 🚫 Blockers

- **<topic>:** <what stops the work, and who or what can unblock it>
```

- **TODO** is what will change, as the user will see it. Don't repeat it in **Solutions**.
- **Solutions** are how to meet the item. If something that already exists meets it, that is the solution; point to it instead of inventing a new one. If the solution is obvious, give the obvious one.
- **Questions** are real decisions the user has to make between two or more options. Never ask the user to approve the TODO or the solution, such as "Go ahead with this?" or "OK with this?", in **Questions** or anywhere else. The user will say if they disagree.
- **Blockers** are things that stop the work, such as a missing API, access or information.
- **Risks** are major problems the change could cause, such as breaking an existing feature or losing data.
- Show **TODO**, **Solutions**, **Questions**, **Blockers** and **Risks** only when they have content.
- The lines under the title count the questions, blockers and risks in this reply. Show a line only for a part that is there.
- Leave out technical details unless the user needs them to decide.
- If something already works, say so and suggest nothing, or only a quick check.
- Answer the user's questions directly, after checking the code.

## Report

When the work is done, reply with only this and nothing else:

```markdown
## <plan title>

✅ <number of issues whose tasks and checks all passed> done\
⏳ <number of issues not done yet, such as skipped or blocked ones> remaining\
❌ <number of issues with a failed task or check> failed

### **#<N>** <title of each issue whose Report section has something in it>

<What that issue's Report asks for.>

### ⏳ Remaining

- **#<N>** <title of each issue not done yet, such as a skipped or blocked one>

### ❌ Failed

- **#<N>** <title of each issue with a failed task or check>: <what failed, and why>
```
