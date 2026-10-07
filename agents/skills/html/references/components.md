# Additional components

Use this reference for details that the form-selection rules in `SKILL.md` do not provide. Read the selected example and `assets/styles.css` when exact markup is unclear.

## Metrics

Use `.metrics` for a small set of supplied measurements that readers must compare or find quickly. Keep each measurement's label, exact value, and unit together. Do not use it as decoration or as a substitute for a table that readers need for exact comparison.

## Markup

Wrap wide tables in `.table-scroll`. Use header cells for row and column headings. Use `<code>` for exact identifiers, paths, commands, and code.

Inline CSS is limited to component values such as `--value`, `--from`, `--to`, and `--at`, and SVG geometry required by a figure. Do not add visual rules that belong in `assets/styles.css`.
