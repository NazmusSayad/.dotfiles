---
name: better-code-full
description: Rethinks the design of existing code, such as data flow, state, and the paths between parts, and rebuilds it in a simpler and more direct way without changing its behavior.

disable-model-invocation: true
metadata: { opencode/autoinvoke: false }
---

## Goal

Find a simpler design for the code in scope and change the code to it. Local tidying, such as inlining a helper or removing a variable, is not the goal. Look for changes that remove whole mechanisms: paths, copies of state, events, layers, and special cases.

Keep the behavior exactly the same. This includes what the user sees, what is stored, what is sent, and edge cases such as timing, ordering, errors, and who receives an update. If a simpler design needs any change in behavior, do not make it. Describe the change and its effect, and ask first.

## Understand the whole design first

Read all the code in scope before you change anything. For a PR, read the full diff and the code it touches. Map:
- Where each piece of state lives, who changes it, and how other parts learn about the change.
- Every path that carries the same data from one part to another.
- Every place that combines, filters, or merges the same data.
- Every special case and the reason it exists.

## Look for a simpler design

Ask how you would build the same behavior from the start with what you know now. Common signs that the design can be simpler:
- The same data reaches a part through more than one path, such as an event, a response field, and a separate list. Use one path.
- The same state is kept in more than one place and kept in sync by hand. Keep one source of truth and compute the rest from it.
- Several places merge or filter the same data in their own way. Do it in one place.
- Separate kinds, events, or types exist for cases that one general rule covers.
- A lookup, query, or subscription exists only to feed another step that could get the value directly.
- Code repeated across implementations, such as providers or handlers, does the same work. Share it when that removes real repetition.
- A layer, wrapper, or indirection only forwards to another.

Before you keep a mechanism, find out why it exists. Code that looks unnecessary often protects behavior, such as an instant open from a cache, an optimistic update, or a provider difference. Keep it when removing it would change behavior.

Prefer one large simplification over many small edits. Fewer lines are good only when the result is easier to follow.

## Make the change

Change the design completely. Remove everything the old design needed: types, events, fields, state, resets, subscriptions, helpers, and their imports. Do not leave the old path beside the new one. Avoid incidental changes to names, formatting, and unrelated code.

## Check the result

Check that each flow in scope behaves as before, step by step, including edge cases. Check that no reference to a removed part remains. Run typecheck and lint.

## Report

- The design before and after, in a few sentences each.
- What the change removed.
- Each part that looks simpler but you kept, and the behavior that it protects.
- The flows that need a manual check.
