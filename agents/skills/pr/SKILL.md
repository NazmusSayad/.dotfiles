---
name: pr
description: "Commits and pushes changes, then creates or updates their pull request."

disable-model-invocation: true
metadata: { opencode/autoinvoke: false }
---

When this run is done, continue normally. Commit or push only when the user asks, and don't remind the user about `pr` or unpublished changes.

Look at the full diff, including uncommitted and new files. Ask if it's unclear which changes belong in the PR.

If you're on the default branch, create a short, descriptive feature branch and switch to it. Otherwise, stay on the current branch.

Commit the PR's changes in the repository's commit style and leave unrelated changes alone. Skip the commit if there is nothing new.

Check for this branch's PR with `gh pr view`. If the PR is merged and the new changes are for the same feature, pull its base branch from origin into this branch and create a new PR against that base branch. If they clearly belong to a different feature, ask what to do.

Push the branch, and ask before force pushing. Then update the open PR, or create one. When this branch has no PR, create it against the default branch.

Read [references/content.md](references/content.md), then write the title and body for the whole PR, not just the latest commit.

If the repository uses a code review bot, request its review if it does not start, and wait for it. Fix each valid point, then commit, push, and update the PR body. Skip a point that is wrong, outside this PR, or against what the user asked for, reply to it on the PR with a short reason, and list it in your final reply. Repeat within this run until no valid points remain.

After success, reply with the status of the commits and the PR, with the PR URL on a new line. Add no other text, fluff, or follow-up messages.
