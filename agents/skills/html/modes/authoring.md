# Authoring

These rules apply to every mode.

## Source

Format existing material. Do not research, inspect a project for findings, or gather new information. Reading a source that the user names is part of formatting it.

Use the file the user names, otherwise material in their message, otherwise the result reached in the current discussion. If none exists, ask for it.

Before writing, list every item the source states. Use the list to place each item and to check coverage at the end.

For a discussion, render where it ended. State settled points as facts, not as a story of the talk. A later message that reverses an earlier one replaces it.

Preserve content that falls outside the mode's usual sections instead of dropping it.

The mode's example shows markup and section order. Its incidental subject matter is not a requirement.

## Fidelity

You may reorder, group, split sentences, and choose the form of each part, but keep any order the source gives. You may not change what the source says. Fidelity wins over every other rule here.

- Add nothing the source does not state.
- Drop nothing that affects the meaning.
- Keep force and certainty. "Must" stays must, and "maybe" stays maybe.
- Keep names, numbers, dates, identifiers, paths, commands, and code exact. Call each thing by one name.
- Label decisions, proposals, assumptions, and open questions as what they are.
- When the source contradicts itself without a resolution, show both statements in a callout.
- Omit what the source never addresses. Where a field is given for some items but not others, write "Not stated". Where the source says TBD, write "Not decided".
- Number items only when the source gives an order.

## Writing

Write every piece of text, including headings, labels, and diagram text, in strict ASD-STE100 Simplified Technical English, with no jargon. There are no exceptions.

Write for the audience and level of detail the user states. If a technical term is unavoidable, explain it in plain words the first time you use it. Give context before detail. Keep related rules and exceptions together.

Every sentence must help the reader understand the document. Delete any sentence that only introduces, repeats, or sums up what is already on the page, and write no captions or footers.

Follow the example's section order and drop the sections the source does not support. Bold only a few key terms. Do not mention the conversation, the source, or this skill.

## Form

Use the simplest form that shows the source's actual structure:

- Reasons and context: paragraphs.
- Parallel items: a bulleted list.
- Ordered actions or phases: `.steps`.
- Three or more items that share two or more fields: a table.
- Events with stated dates: `.timeline`.
- Work with stated start and end dates: `.gantt`.
- Stated branches, loops, parallel paths, or links among four or more parts: an SVG diagram.
- Given numbers worth comparing: `.bars`.
- Conflicts, blockers, or open points inside a section: `.callout`.
- Optional evidence: `<details>`.

A linear sequence is a list, not a diagram. Undated phases never become a timeline or Gantt chart. A diagram answers one question, has at most about ten nodes, a verb on each arrow, and no node or arrow the source does not state. The document reads completely without its figures.

Use Font Awesome icons where the example uses them. Keep the contents list only when there are four or more sections.

Use `.metrics` for a small set of supplied measurements that readers must compare or find quickly. Keep each measurement's label, exact value, and unit together. Do not use it as decoration or as a substitute for a table that readers need for exact comparison.

## Reply

Also name gaps, unresolved conflicts, structure you inferred, and content you left out.
