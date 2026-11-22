---
name: better-code
description: Rules for writing simple, direct, and correct code, including TypeScript and React. MUST USE for writing, editing, fixing, reviewing, or refactoring code in any language.
---

## Writing code

Follow the project's existing patterns. Give each piece of state one owner and derive the rest. Fix the data shape instead of piling up flags and branches.

Use the simplest clean, readable, and direct solution that fully satisfies the request. Prioritize simplicity, readability, and directness over reusability, abstraction, and all other forms of overengineering. Follow YAGNI principles: do not add features, flexibility, or complexity before they are required. Do not introduce variables, functions, helpers, interfaces, types, or other abstractions unless they simplify complex logic or remove substantial repetition.

Use explicit logic to prevent ambiguity and implicit fallbacks. For example, avoid `if "yes": true; else: false`; use `if "yes": true; if "no": false; else: exception`. Do not hide errors with fallbacks, retries, or empty catches. Handle errors explicitly and fail fast.

## Changes

Keep changes focused. Change only what is necessary to satisfy the user's request and produce a correct result.

Do not refactor, clean up, reformat, rename, or otherwise improve unrelated work. Remove code made unused by your changes, but don't remove pre-existing dead code unless asked.

## Comments and Tests

Do not write comments unless the code does something it cannot explain by itself.

Do not add tests (eg: unit, integration, regression, or smoke tests) to the project unless explicitly requested.

## TypeScript

Before working with TypeScript (`.ts`, `.tsx`, `.mts`), read [references/typescript.md](references/typescript.md).

## React

Before working with React code, components, hooks, or JSX (`.jsx`, `.tsx`), read [references/react.md](references/react.md).
