---
name: uncomment
description: Removes code comments, keeping only those without which a reader would likely break the code.

disable-model-invocation: true
metadata: { opencode/autoinvoke: false }
---

Delete comments, not code. Do not rename, restructure, or rewrite code to make up for a removed comment.

Judge each comment on its own. Keep it only when, without it, a reader would likely break the code or "fix" something that is intentional, such as a platform quirk, a hidden constraint, or a security guard. Everything else goes, including comments that explain intent or describe what the code does.

Keep comments that tools read, such as linter and type-checker directives, pragmas, and license headers, because removing them changes behavior.
