---
name: variants
description: Creates several unique variants of one idea, one per subagent.

disable-model-invocation: true
metadata: { opencode/autoinvoke: false }
---

A variant is one unique take on an idea, created by its own subagent. Each one must look good and bring value to the user, even when all of them are completely different. The user may give some references; a reference is an idea, not something to copy. Create the number of variants the user asks for, defaulting to 3.

## Subagents

Use one subagent per variant. Run the subagents in the background.

Tell each subagent what the idea is for, what already exists, what the user rejected, and only the requirements the result truly needs. Then tell it to be creative. Do not explain how to make its variant, give it your bias, or push it toward one design. Do not add "you must" or "you cannot" rules about creative choices.

Guide a subagent only when it goes too far: when it does something wrong or rewrites everything.

## Uniqueness

Variants that look the same, or differ only in small details, have failed.

When the user deletes or rejects variants, remove them and do not bring them back. When the user asks for more, make the new ones more creative than the last.
