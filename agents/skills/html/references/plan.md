# Plan profile

## Source state

Before writing, list every item the plan states: problem and background, outcome, scope, restrictions, steps, decisions and reasons, rejected options, dependencies, owners, dates, statuses, assumptions, risks, open questions, files, commands, and completion checks. Use the list to place each item and to check coverage at the end.

For a discussion, render where it ended. State settled points as facts, not as a story of the discussion. A later statement that reverses an earlier statement replaces it.

## Fidelity

You may reorder, group, split sentences, and choose the form of each part, but keep any order the plan gives. Do not change what the plan says. Fidelity wins over every other rule in this profile.

- Add nothing the plan does not state: no task, owner, date, estimate, status, priority, dependency, risk, or metric.
- Drop nothing that affects the work.
- Leave out features and tasks that are not part of the plan, including items the plan lists as out of scope or not needed. Show something as not allowed only when the plan forbids it directly and explicitly.
- Keep force and certainty. “Must” stays must, “maybe” stays maybe, and planned work never reads as done.
- Keep names, numbers, dates, identifiers, paths, commands, and code exact. Call each thing by one name.
- Label decisions, proposals, assumptions, and open questions as what they are.
- When the plan contradicts itself without a resolution, show both statements in a callout.
- Omit what the plan never addresses. Where a field is given for some items but not others, write “Not stated.” Where the plan says TBD, write “Not decided.”
- Number items only when the plan gives an order.

## Writing

Follow the example's section order and drop sections the plan does not support. Do not mention the discussion or the source.

## Form

Use the simplest form that shows the plan's actual structure:

| Plan content | Form |
|---|---|
| The problem the plan solves | Problem section: paragraphs for the background, and charts or diagrams when the plan gives numbers or structure |
| Reasons and context | Paragraphs |
| Parallel items | Bulleted list |
| Ordered actions or phases | `.steps`; each step says what it changes and, when the plan states them, its rules, dependencies, completion check, and undo |
| Three or more items that share two or more fields | Table |
| The parts of the plan and their results | Summary table; add a status column only when the plan gives statuses |
| Files to change | Folder tree with one very short line per file; omit it when the list does not help, such as for very large changes |
| Completion checks | `.checklist` |
| Events with stated dates | `.timeline` |
| Work with stated start and end dates | `.gantt` |
| Stated branches, loops, parallel paths, or links among four or more parts | SVG diagram |
| Given numbers worth comparing | `.bars` |
| Decisions, rejected options, risks, assumptions, and open questions | One bulleted list in the Notes section, each item opening with the item in bold and an icon that shows its kind |
| Conflicts, blockers, or open points inside a section | `.callout` |
| Optional evidence | `<details>` |

A linear sequence is a list, not a diagram. Undated phases never become a timeline or Gantt chart. A diagram answers one question, has at most about ten nodes, has a verb on each arrow, and contains no node or arrow that the plan does not state. The document must read completely without its figures.

## Template

Use `examples/plan.html` as the template and replace its content with the plan. Keep its two stylesheet links and its script, and add no other external assets. For components the example does not show, such as `.timeline`, `.gantt`, `.bars`, `.callout`, `.status`, `.diagram`, and `<details>`, read their rules in `references/components.md` and `assets/styles.css` before writing the markup.

Use Font Awesome icons where the example uses them. Choose each icon for its meaning and put a `text-*` color class on the icon itself. Section headings have no icons. In the contents list, use outline icons (`fa-regular`). The free outline set is small, so use only icons it includes.

Keep the contents list only when there are four or more sections. Give sections and steps `id` values so cross-references become links. In diagrams, fit the `viewBox` to the drawing, keep 40 or more units between nodes, and keep labels clear of lines.
