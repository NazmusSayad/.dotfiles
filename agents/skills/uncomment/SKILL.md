---
name: uncomment
description: Removes code comments, keeping only those without which the code cannot be understood.

disable-model-invocation: true
metadata: { opencode/autoinvoke: false }
---

Delete comments, not code. Do not rename, restructure, or rewrite code to make up for a removed comment.

Keep a comment only when a reader could not understand the code without it. A comment that is merely helpful, explains intent, or describes what the code does goes. When in doubt, remove it.

Keep comments that tools read, such as linter and type-checker directives, pragmas, and license headers, because removing them changes behavior.
