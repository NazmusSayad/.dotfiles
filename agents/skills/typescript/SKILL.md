---
name: typescript
description: TypeScript best practices for clean, maintainable, and optimized code. MUST USE for writing or working with TypeScript code (.ts, .tsx, .mts files), including editing, reviewing or refactoring.
---

## Types

Avoid explicit type annotations when TypeScript can infer.
Do not use `any`, casts, or explicit generic type arguments when inference is sufficient.

## Variables

Use consistent, descriptive naming; avoid obscure abbreviations.
Do not use object or array destructuring in declarations, assignments, parameters, or loops; use direct property or indexed access instead, including for component props and tuple returns.

## Asynchronous

Prefer `async`/`await` over callbacks or `.then()` chains
