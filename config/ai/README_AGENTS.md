## Code Style and Best Practices

- Keep code minimal, direct, and readable.
- Do not add comments unless explicitly instructed.
- Do not use hacks, workarounds, magic values, or undocumented behavior.
- Implement logic in the most straightforward and explicit way supported by the language or framework.

### Abstractions

- Prefer direct, inline implementations.
- Avoid introducing helpers, wrappers, or abstractions unless they clearly simplify the code.
- Only introduce abstractions when they remove significant duplication or substantially improve readability.

### TypeScript

- Prefer type inference whenever possible.
- Do not use `any`, casts, or explicit generic type arguments when inference is sufficient.

## CLI / Shell / Terminal

- Use `bash` for all shell commands.
- Use Unix-style paths when referencing files or directories.
- Do not write scripts to modify file contents. Use proper tools instead.
