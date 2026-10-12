---
name: pr
description: "Commits and pushes changes, then creates or updates their pull request."

disable-model-invocation: true
metadata: { opencode/autoinvoke: false }
---

When this run is done, continue normally. Commit or push only when the user asks. Don't offer or remind the user to commit, push, or run `pr`, and don't mention unpublished changes.

Look at the full diff, including uncommitted and new files. Ask if it's unclear which changes belong in the PR.

If you're on the default branch, create a short, descriptive feature branch and switch to it. If you're on a feature branch and the changes are clearly for a different feature, create a new feature branch from the base branch and switch to it. If you're not sure, ask. Otherwise, stay on the current branch.

Commit the PR's changes in the repository's commit style and leave unrelated changes alone. Skip the commit if there is nothing new.

Check for this branch's PR with `gh pr view`. If the PR is merged and the new changes are for the same feature, pull its base branch from origin into this branch and create a new PR against that base branch.

Push the branch, and ask before force pushing. Then update the open PR, or create one. When this branch has no PR, create it against the default branch.

Read [references/content.md](references/content.md), then write the title and body for the whole PR, not just the latest commit.

After success, reply with the status of the commits and the PR, with the PR URL on a new line. Add no other text, fluff, or follow-up messages.
