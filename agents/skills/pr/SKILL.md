---
name: pr
description: "Commits and pushes changes, then creates or updates their pull request."

disable-model-invocation: true
metadata: { opencode/autoinvoke: false }
---

Each run publishes once. After that, keep working normally and don't commit or push again until `pr` is run again.

## Commit and push

Look at the full diff, including uncommitted and new files. Ask if it's unclear which changes belong in the PR.

If you're on the default branch, create a short, descriptive feature branch and switch to it. Otherwise, stay on the current branch.

Commit the PR's changes in the repository's commit style and leave unrelated changes alone. Skip the commit if there is nothing new. Push the branch, and ask before force pushing.

## Pull request

Check for this branch's PR with `gh pr view`. If it's open, update it. If there is none, create one against the default branch. If it was merged or closed, ask what to do.

Read [references/content.md](references/content.md), then write the title and body for the whole PR, not just the latest commit.

Return the PR link and any problem.
