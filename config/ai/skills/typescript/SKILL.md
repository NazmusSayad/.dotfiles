---
name: typescript
description: TypeScript code style and optimization guidelines. MUST USE before writing or modifying any TypeScript code (.ts, .tsx, .mts files). Also use when reviewing code quality or implementing type-safe patterns. Triggers on any TypeScript file edit, code style discussions, or type safety questions.
---

## Best Practices

- Prefer type inference whenever possible.
- Do not use `any`, casts, or explicit generic type arguments when inference is sufficient.
- Prefer `async`/`await` over callbacks or `.then()` chains
