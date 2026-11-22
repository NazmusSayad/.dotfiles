## Types

Avoid explicit type annotations when TypeScript can infer.
Do not use `any`, casts, or explicit generic type arguments when inference is sufficient.

## Variables

Use consistent, descriptive naming; avoid obscure abbreviations.
Do not use object or array destructuring in declarations, assignments, parameters, or loops; use direct property or indexed access instead, including for component props and tuple returns.

## Asynchronous

Prefer `async`/`await` over callbacks or `.then()` chains.
