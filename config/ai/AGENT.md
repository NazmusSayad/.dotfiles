## Files

- If no file is specified, prioritize the currently opened file.
- Modify existing files only. Do not create new files or directories unless instructed.

## Code Style

- Keep code minimal, direct, and readable.
- Do not add comments unless explicitly instructed.
- Do not use hacks, workarounds, magic values, or undocumented behaviors. Implement logic in the most straightforward, explicit way supported by the language or framework.

### Abstractions

- Prefer the most direct, inline implementation.
- Do not introduce abstractions, helper functions, extra variables, or additional structures unless they clearly improve readability or correctness in idiomatic TypeScript. Allow small helpers only when they reduce duplication or clarify non-trivial logic.

## TypeScript Rules

- Rely on full type inference.
- Do not use type `any`, casts, or explicit generic type arguments when inference is possible.

## npm packages

Use @antfu/ni to manage npm packages.

### `nr` - run

```bash
nr dev --port=3000
nr build
```

### `ni` - install

```bash
ni <package>
ni -D <package>
```

### `nup` - upgrade

```bash
nup
nup <package>
```

### `nun` - uninstall

```bash
nun <package>
```

## CLI/Terminal

- Run terminal commands only when instructed.
- Never prefix terminal commands with `cd`. Use `cd` only when entering a subdirectory.

### Example

Incorrect:

```sh
cd ./project && nr test
```

Correct:

```sh
nr test
```

Correct:

```sh
cd ./project/subproject && nr test
```
