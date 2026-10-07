# Components

Use components for the relationship they express, not for decoration. Read the selected example and `assets/styles.css` when exact markup is unclear.

## Selection

| Content relationship | Component |
|---|---|
| A few document facts | `.brief` or `.meta` |
| A few central measurements | `.metrics` |
| Parallel items | Bulleted list |
| Ordered actions | `.steps` |
| Completion conditions | `.checklist` |
| Three or more items sharing two or more fields | Table inside `.table-scroll` |
| Events with stated dates | `.timeline` |
| Work with stated start and end dates | `.gantt` |
| Comparable supplied numbers | `.bars` |
| Branches, loops, or relationships among several parts | SVG `.diagram` |
| A conflict, limit, warning, or important exception | `.callout` |
| A supplied status | `.status` |
| Useful but optional evidence or detail | `<details>` |

Do not use a timeline or Gantt chart for undated material. A linear sequence is normally a list. Do not place prose in a table unless comparison across columns helps the reader.

## Visual truth

Figures must represent only relationships and values stated by the source. Use a zero baseline for bars unless the source explicitly requires another baseline and the document identifies it. Do not compare values with different units on one scale.

A diagram answers one question, contains at most about ten nodes, and uses a verb on each arrow. Fit its `viewBox` to the drawing, keep at least 40 units between nodes, and keep labels clear of lines. Include the same facts in nearby prose or a list.

## Markup

Use a `<table>` with header cells for tabular data. Give times a `<time>` element when their machine-readable value is known. Use `<code>` for exact identifiers, paths, commands, and code. Use `<strong>` sparingly for text readers must locate quickly.

Use icons only to reinforce text. Do not use color or an icon as the only status, warning, or category indicator. Keep link text descriptive. Provide accessible text for every exact value shown visually.

Keep the layout usable at narrow widths and in print. Do not add inline CSS except component custom properties such as `--value`, `--from`, `--to`, and `--at`, and SVG geometry required by a figure.
