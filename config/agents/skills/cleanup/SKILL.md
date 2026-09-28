---
name: cleanup
description: Removes unnecessary code and content, simplifies what remains, and clears out the leftovers without losing required behavior or details.

disable-model-invocation: true
metadata: { opencode/autoinvoke: false }
---

## Remove unnecessary parts

Trace references and dependencies before deciding code is unused. Public APIs, routes, and configuration can expose code without local callers. Preserve stored data and schemas needed by the behavior that remains.

Delete unnecessary parts instead of rebuilding them behind new layers. When removing a feature, follow its dependencies: remove unused imports, exports, variables, helpers, files, assets, and configuration. Keep dependencies still used elsewhere. Remove leftover aliases, shims, placeholders, and commented-out code that no longer serve a required purpose.

## Simplify the remaining logic

Prefer direct, explicit logic. Remove one-use wrappers and helpers when inlining them reads better; keep shared code when inlining would create substantial repetition. Introduce abstractions only when they simplify complex logic or remove substantial repetition. Replace comments that merely explain convoluted code with clearer code.

Avoid incidental changes to names, formatting, schemas, or unrelated implementations. They add work to review without simplifying the target. Fewer lines, files, or layers help only when the result is easier to read.

## Organize only where it helps

Separate distinct responsibilities where that makes the code easier to follow, rather than separating every function or component. Keep related small helpers and subcomponents beside the code they serve. Extract a file when its responsibility warrants the separation, not merely to make code reusable. Splitting one large file into many tiny files can make it harder to follow.

Separate data construction or transformation from rendering when mixing them obscures the logic. Pages and entry points are often clearer when focused on wiring the relevant parts together; small sibling helpers can stay there. Sharing code is worthwhile when it removes substantial repetition, not when it requires changing another implementation solely to reuse a small piece.

## Clean up content

Preserve facts, instructions, examples, links, routes, and details the audience needs. Cut repetition, false claims, unnecessary implementation detail, and filler. Link to existing documentation when it already explains a supporting topic, while keeping enough context for the current explanation to make sense. Check claims about features, commands, benchmarks, and product behavior against their sources. Keep examples accurate rather than inventing them to fill gaps.

## Check the result

Check that removed features have no remaining entry points or orphaned dependencies, no live references point to removed symbols, and retained behavior and content still meet their requirements.
