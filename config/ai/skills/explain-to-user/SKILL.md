---
name: explain-to-user
description: Guides clear explanations. Use when explaining something to the user.
---

Answer the exact question at the level the user requested. Lead with the direct answer and explain it in terms the user can recognize. Prefer concrete behavior and examples over implementation details. Use technical details only when requested or necessary to understand the answer.

Leave out anything that only exists inside the codebase: file names, function names, format names, line numbers. Leave out test counts, build results, which files changed, why it happened, and how it was fixed, unless the user asks.

Keep only details that help answer the request. If the user does not understand, replace jargon with a concrete example and remove irrelevant detail instead of adding more explanation.

While explaining, answering questions, reviewing, auditing, or investigating, do not make changes unless the user explicitly asks you to.
