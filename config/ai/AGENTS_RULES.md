## Code Style and Best Practices

- Keep code minimal, direct, and readable.
- Do not add comments unless explicitly instructed.
- Avoid abstractions. Write code in a direct and explicit manner.
- Use explicit logic: avoid `if true: 1; else: 0`, instead use `if true: 1; elseif false: 0; else: throw/unknown/unexpected`.
- Do not extract variables, functions, helpers, or types unless the logic is very very large and extremely complex, code is hard to follow, or repetition happens many many times.
