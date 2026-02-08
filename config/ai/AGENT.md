# Files

- If no file is specified, prioritize the currently opened file.
- Modify existing files only. Do not create new files or directories unless instructed.

# Code Style and Best Practices

- Keep code minimal, direct, and readable.
- Do not add comments unless explicitly instructed.
- Do not use hacks, workarounds, magic values, or undocumented behaviors.
- Implement logic in the most straightforward, explicit way supported by the language or framework.

## Avoid Abstractions and Over Declarations

- Prefer the most direct, inline implementation.
- Add helpers, extra variables, or abstractions only when they reduce meaningful duplication.

## TypeScript Rules

- Rely on full type inference.
- Do not use type `any`, casts, or explicit generic type arguments when inference is possible.

# Skills

- Use `npm` skill to work with npm.
- Use `react` skill when writing React code.
- Use `frontend-design` skill when designing frontend interfaces.

# CLI/Terminal

- Run terminal commands only when instructed.
- Never prefix terminal commands with `cd`. Use `cd` only when entering a subdirectory.

## Example

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
