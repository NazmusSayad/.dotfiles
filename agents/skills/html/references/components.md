# Components

`assets/styles.css` defines every component. The mode examples in `modes/` show most of them. For a component they do not show, read its rules in the stylesheet to build the markup. Take only markup from the examples, never their content or section order.

- `.doc-header` with `h1` and `.lede`: a page title and an opening line.
- `.contents`: links to the page's sections.
- `main > section` with `h2`: a top-level section.
- `.message` with `.message-meta`: one message, comment, or post.
- `.brief`, `.meta`, `.fields`: labels with values.
- `.metrics`: headline measurements.
- `.steps`: ordered actions or phases.
- `.checklist`, `li.done`: checks, open or done.
- `.timeline`: dated events.
- `.gantt`: work with start and end dates.
- `.bars`: numbers to compare.
- `.table-scroll`: a wrapper for a wide table.
- `.callout`, `.callout.warn`: a point set apart from the text.
- `.status.ok`, `.status.warn`, `.status.risk`: a status label.
- `.figure-frame` with `svg.diagram`: an SVG diagram.
- `<details>` with `<summary>`: a collapsible block.

## Message

```html
<article class="message">
  <header class="message-meta">Author, time, or link, as the source shows them</header>
  <p>Body</p>
</article>
```

Put the body's elements directly inside the `article`. Leave out `.message-meta` when the source shows no author, time, or link.

## Icons and links

Color each Font Awesome icon with a `text-*` class on the icon itself. Section headings have no icons. In the contents list, use outline icons (`fa-regular`). The free outline set is small, so use only icons it includes. Give sections and steps `id`s so cross-references become links. In diagrams, fit the `viewBox` to the drawing, keep 40 or more units between nodes, and keep labels clear of lines.

## Markup

Wrap wide tables in `.table-scroll`. Use header cells for row and column headings. Use `<code>` for exact identifiers, paths, commands, and code.

Inline CSS is limited to component values such as `--value`, `--from`, `--to`, and `--at`, and SVG geometry required by a figure. Do not add visual rules that belong in `assets/styles.css`.
