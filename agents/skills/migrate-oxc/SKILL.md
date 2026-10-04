---
name: migrate-oxc
description: Migrate ESLint and Prettier projects to Oxlint and Oxfmt.

disable-model-invocation: true
metadata: { opencode/autoinvoke: false }
---

# Migrate to Oxc

Use [create-src](https://github.com/NazmusSayad/create-src) as the reference. Inspect its current root and template configurations rather than copying details from this skill.

## Goal

Replace ESLint and Prettier completely while preserving supported rules, severities, options, ignores, overrides, and project integrations. Prefer the simplest native Oxc setup that does so.

## Migration

Inspect the whole repository before editing. Include generated templates, package scripts, hooks, CI, editor settings, generator code, documentation, dependencies, and committed lockfiles.

Follow the reference architecture:

- Keep `oxlint.config` and `eslint-recommended` together at each project root.
- Use `.mts` for ESM-style TypeScript configs in CommonJS packages; otherwise follow the package module format.
- Keep an explicit ESLint-compatible recommended rule map where Oxlint's preset is not equivalent.
- Prefer native Oxlint rules and plugins. Use JS plugins only for gaps the project actually relies on.
- Do not build custom compatibility layers for unsupported ESLint rules. Preserve what Oxc supports and clearly identify intentional losses.
- Keep each generated template self-contained.
- Keep `.oxfmtrc.json` minimal and project-specific.

For type-aware linting, use `oxlint-tsgolint` and `options.typeAware`. Resolve reported TypeScript configuration incompatibilities rather than disabling type-aware linting. Remember that nested template configs may need validation with an explicit config path, such as `oxlint -c ./oxlint.config.mts .`.

## Commands and integrations

Use `eslint-plugin-oxfmt` so Oxlint checks and fixes JS/TS formatting:

- `lint` runs type checking when the project already requires it, then Oxlint.
- `lint:fix` runs only `oxlint --fix`.
- Do not add `lint:format` or follow Oxlint with Oxfmt for code files.
- In lint-staged, use Oxlint for code and Oxfmt only for non-code formats such as JSON, Markdown, YAML, and CSS.
- Generated-project finalization should likewise avoid a second code-formatting pass.

Remove obsolete direct ESLint and Prettier tooling, then add only the Oxc packages and JS plugins used by the resulting configuration. ESLint may remain transitively through a JS plugin. Update existing lockfiles and relevant editor recommendations.

## Verification

Run the final, user-facing workflows after all edits:

- Root and template `lint` and `lint:fix` commands
- Type checking where applicable
- Template configs in their standalone context
- Oxfmt checks for supported non-code files
- Hook and generator commands

Search once more for stale ESLint or Prettier configs, dependencies, scripts, documentation, `lint:format`, and duplicate `oxlint`/`oxfmt` pipelines.

Do not claim a hook was tested end to end unless it was run with staged files. Otherwise, state that its underlying commands were verified.
