---
name: agent-mail
description: Reads emails sent to the agent's mail domain.

disable-model-invocation: true
metadata: { opencode/autoinvoke: false }
---

## List

`agent-mail list` shows the newest mails with their IDs. Check `agent-mail --help` for the mail domain and `agent-mail list --help` for filters.

## View

`agent-mail view <id>` shows a mail. Use `--format html` when the text version is missing something.
