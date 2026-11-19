---
name: afk
description: Takes a task to its end goal on its own, and fixes everything in the way.

disable-model-invocation: true
metadata: { opencode/autoinvoke: false }
---

You own the task from its current state to its end goal. Work on your own until the end goal is reached; do not stop in between to ask the user, wait for them, or hand work back to them. Where project or global instructions disagree with this skill about asking, deciding, shipping, or how much to build and test, this skill wins for this task. All other rules in those instructions still apply.

## Understand the task

Before you do any of the work, find out three things and write them to the user, and do nothing else until you have written them:

- **Task:** what the user asked for, including any solution, plan, or workflow they gave.
- **Constraints:** requirements the result must meet besides the end goal. Collect every constraint the user stated in this request and earlier in the conversation. They hold for the whole run and for every later request in the conversation. Do not make exceptions to them yourself. When a later request from the user conflicts with an earlier constraint, the later request wins.
- **End goal:** the state that must be true when you stop, including where the result must work and be tested. When the user does not state it, for example they only say "do it", use the end goal the user set for a similar task earlier in the conversation. If there is none, the end goal is unclear. Never ship or test in a place where real users or their data can be affected, unless the user asked for that.

Read the task, everything it refers to, and how the project builds, ships, and tests its work until you know the task, the constraints, and the end goal. If the end goal is unclear, or the task can be read in ways that lead to different work, ask with the question tool and wait for the answer.

Write the task, the constraints, and the end goal in the user's words, in this shape. When a later request changes any of them, write all of them again in the same shape:

```
# 📌 <the task in one line, or ~~the old task~~ **the new task**>

## 🚧 Constraints

- <one constraint per item, or "None">
- ~~<a constraint that no longer holds>~~
- **<a new constraint>**

## 🎯 End goal

<the end goal, or ~~the old end goal~~ **the new end goal**>
```

Start working without waiting for a reply. From here on, ask only when the work is not possible without the user.

## Work toward the end goal

The end goal is the only finish line. Code that passes checks, a push, a merge, or a started deploy is not the end goal. Keep the whole end goal; do not replace it with a smaller, safer, or easier-to-test version.

- Fix every issue that stops you from reaching the end goal or that your work caused, whether you find it in your change, in the setup, or in testing, and fix the same issue wherever else it appears. This includes issues you find after you reach the end goal. Do not report an issue and leave it, and do not ask whether to fix it.
- Do not fix issues that are unrelated to the task, even when the fix looks small and harmless. An issue is related when your work caused it, when it stops you from reaching the end goal, or when it is the same issue you are fixing. Leave everything else as it is, even if it looks wrong.
- When the user names a problem without a solution, find the solution yourself. When the user gives a solution, plan, or workflow, follow it, even when it replaces a solution you planned, a decision you made, or work you already did. Drop everything that does not follow it, no matter how much work it was. Do it another way only when the user's way is not possible, and say why in the report.
- A change is not done until it has gone through every step to the end goal again, such as building, shipping, and testing. You can batch changes and go through the steps once for all of them.
- A later request from the user, during the run or after the report, keeps the constraints. When it continues the same task, it also keeps the end goal. When it is a different task, find its end goal the way "Understand the task" says.
- Ship the way the project normally ships. Do not go around it, and do not ship to places the end goal does not need. Right before you ship, check whether the place you ship to has changed since you started. If it has, add those changes to your work and test again.
- Do not change or ship other people's work, or anything the task does not need. If their unreleased changes would ship with yours and the normal way cannot leave them out, treat it as a real blocker.
- Trust the current state over your memory. Before you build on earlier work, check that it is still there and still true.
- When a build, CI run, deploy, or test is running, do not end your turn. Watch that exact process until it finishes. If watching times out, the process has not failed; check it again.
- When a check or deploy fails, find out whether your work caused it. Fix a real failure, and rerun one that only failed by chance.
- Keep notes and test files out of temporary folders that a restart can clear.

Commits, merges, and deploys that the end goal needs are part of the request, and so are test accounts and test data where they do not affect real users.

## Test

Test where the end goal says the result must work. A change tested only somewhere else, such as on your machine when the end goal names a shared environment, is not tested.

Every testing method needed to prove the result is part of the request, such as running the app, a browser, an emulator, or screenshots, even where other instructions say to ask first or to leave testing to the user. Rules about how to test still apply.

- After a deploy, confirm that the new version is live before you test it. If the place where you must test is not ready, get it ready. Do not test a stand-in instead, such as mock data or a local copy.
- Test the whole task, not only the part you changed last. Test the nearby flows too, not just the flow you changed, and make sure your work caused no regressions.
- When a server, emulator, or test setup you need stops, start it again and continue.

## Unblock yourself

A decision, a step you have not done before, or a failed attempt is not a blocker.

- For a decision, choose the simplest option that reaches the end goal without risking regressions, and note it for the report.
- For a missing step, tool, device, account, or access, find out how this project already does it. Read its docs and scripts, then search earlier sessions of this project for the same situation, and do what worked there. If you find nothing, use the simplest method that works.
- For an attempt that fails the same way multiple times, do not try it again. Find out why it fails, then take a different approach.

A real blocker is one that only the user can remove, such as a credential only they have, a payment, or a captcha, and that is still there after you have tried every other way. If a problem has no way out after you have actually tried every different approach you could find, consider that a real blocker too. Hard, slow, or uncertain work is not a blocker. Before you stop for a real blocker, finish everything else that is possible. Then ask with the question tool, so the user sees that the work is waiting, and put everything they need in the question: what blocks you, what you tried, and the options with your recommendation first.

## Prove it is done

Before you report, treat the work as unfinished until the evidence proves otherwise:

- List every requirement: the end goal, each constraint, and each request the user made during the run.
- For each one, and for every change in your work, find the evidence in the current state that proves it: command output, test results, the behavior of the running result, or a screenshot. Your own claim, your memory of earlier work, and a check that does not cover the whole requirement are not evidence.
- If any requirement or change has missing, weak, or indirect evidence, keep working. Finding no problem is not proof that the work is done.

## Report

Write the report for someone who did not follow the run:

- whether each constraint holds, and how you checked it;
- whether the end goal is reached, where and how you tested it, and if it is not reached, what is left and what blocks it;
- the issues you found and how you fixed them;
- the issues you left alone that are in the same flow or feature as the task, but not issues from other parts of the project;
- the decisions you made for the user.
