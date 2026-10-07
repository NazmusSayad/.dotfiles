# Components

## Form

Use the simplest form that shows the source's actual structure:

| Source content | Form |
|---|---|
| Reasons and context | Paragraphs |
| Parallel items | Bulleted list |
| Ordered actions or phases | `.steps` |
| Three or more items that share two or more fields | Table |
| Events with stated dates | `.timeline` |
| Work with stated start and end dates | `.gantt` |
| Stated branches, loops, parallel paths, or links among four or more parts | SVG diagram |
| Given numbers worth comparing | `.bars` |
| Conflicts, blockers, or open points inside a section | `.callout` |
| Optional evidence | `<details>` |

A linear sequence is a list, not a diagram. Undated phases never become a timeline or Gantt chart. A diagram answers one question, has at most about ten nodes, a verb on each arrow, and no node or arrow the source does not state. The document reads completely without its figures.

For components the example does not show, such as `.timeline`, `.gantt`, `.bars`, `.callout`, `.status`, `.diagram`, and `<details>`, read their rules in `assets/styles.css` to build the markup.

## Icons, contents, and links

Use Font Awesome icons where the example uses them, each chosen for its meaning and colored with a `text-*` class on the icon itself. Section headings have no icons. In the contents list, use outline icons (`fa-regular`). The free outline set is small, so use only icons it includes. Keep the contents list only when there are four or more sections. Give sections and steps `id`s so cross-references become links. In diagrams, fit the `viewBox` to the drawing, keep 40 or more units between nodes, and keep labels clear of lines.

## Metrics

Use `.metrics` for a small set of supplied measurements that readers must compare or find quickly. Keep each measurement's label, exact value, and unit together. Do not use it as decoration or as a substitute for a table that readers need for exact comparison.

## Markup

Wrap wide tables in `.table-scroll`. Use header cells for row and column headings. Use `<code>` for exact identifiers, paths, commands, and code.

Inline CSS is limited to component values such as `--value`, `--from`, `--to`, and `--at`, and SVG geometry required by a figure. Do not add visual rules that belong in `assets/styles.css`.
