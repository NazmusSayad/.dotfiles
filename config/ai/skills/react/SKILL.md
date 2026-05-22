---
name: react
description: React code style, conventions, best practices, and patterns for building clean, efficient, and scalable applications. MUST USE for writing or working with any React code, components, hooks, or JSX (.jsx, .tsx files), including editing, reviewing or refactoring.
---

## Hooks

- Avoid `useEffect` for simple state updates or derived state. Use `useMemo`, `useCallback`, or direct calculations instead.

## Minimize Over-Declaration

Inline simple logic. Extract only when truly big and complex.

❌ Over-declared:

```tsx
function MyComponent() {
  const [count, setCount] = useState(0)

  function handleIncrement() {
    setCount(count + 1)
  }

  return (
    <div>
      <p>Count: {count}</p>
      <button onClick={handleIncrement}>Increment</button>
    </div>
  )
}
```

✅ Inlined:

```tsx
function MyComponent() {
  const [count, setCount] = useState(0)

  return (
    <div>
      <p>Count: {count}</p>
      <button onClick={() => setCount(count + 1)}>Increment</button>
    </div>
  )
}
```

## Pure Functions Outside Components

Pure utility functions must be declared outside React components.

❌ Wrong:

```tsx
function MyComponent() {
  const [text, setText] = useState("")

  function processText(text: string) {
    return text.toUpperCase()
  }

  return <p>{processText(text)}</p>
}
```

✅ Right:

```tsx
function MyComponent() {
  const [text, setText] = useState("")

  return <p>{processText(text)}</p>
}

function processText(text: string) {
  return text.toUpperCase()
}
```

## Tailwind Classes

Prefer direct inline Tailwind classes in JSX. Do not extract them into variables or helpers. Only extract when the same classes are reused extensively or the logic becomes very complex.

❌ Unnecessary extraction:

```tsx
function MyComponent() {
  const containerClasses = "flex flex-col gap-4 p-6"
  return <div className={containerClasses}>...</div>
}
```

✅ Inline:

```tsx
function MyComponent() {
  return <div className="flex flex-col gap-4 p-6">...</div>
}
```

✅ Complex conditionals, use `cn` utility:

```tsx
import { cn } from "@/lib/utils"

function Form({ hasError }: { hasError: boolean }) {
  return (
    <input
      className={cn(
        "rounded-lg border-2 px-4 py-2 transition-colors",
        hasError
          ? "border-red-500 bg-red-50 focus:ring-red-200"
          : "border-gray-300 focus:ring-blue-200"
      )}
    />
  )
}
```
