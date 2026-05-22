---
name: typescript
description: TypeScript best practices for clean, maintainable, and optimized code. MUST USE for writing or working with TypeScript code (.ts, .tsx, .mts files), including editing, reviewing or refactoring.
---

## Types

- Prefer type inference whenever possible.
- Do not use `any`, casts, or explicit generic type arguments when inference is sufficient.

## Asynchronous

- Prefer `async`/`await` over callbacks or `.then()` chains
