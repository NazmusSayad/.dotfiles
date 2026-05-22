---
name: react
description: React code style, conventions, best practices, and patterns for building clean, efficient, and scalable applications. MUST USE for writing or working with any React code, components, hooks, or JSX (.jsx, .tsx files), including editing, reviewing or refactoring.
---

## Hooks and state

- Do not use `useEffect` for simple state updates or derived state.
- Prefer direct calculations, `useMemo`, or `useCallback` when appropriate.
- Keep effects for real side effects: subscriptions, network sync, DOM/browser APIs, timers, external systems.

## Avoid unnecessary declarations

- Inline simple logic directly in JSX.
- Do not extract simple handlers, variables, components, or helpers just for naming.
- Declare functions only when logic is complex, reused, or improves readability.

```tsx
function Counter() {
  const [count, setCount] = useState(0)

  return <button onClick={() => setCount(count + 1)}>Count: {count}</button>
}
```

Use a named handler only when the body is non-trivial:

```tsx
function Counter() {
  const [count, setCount] = useState(0)

  function handleIncrement() {
    // complex logic here eg: analytics, validation, side effects, 50 lines of code, etc.
    setCount(count + 1)
  }

  return <button onClick={handleIncrement}>Count: {count}</button>
}
```

## Keep pure utilities outside components

- Pure functions that do not depend on component scope must live outside the component.
- This avoids recreating them on every render and keeps components focused.

```tsx
function Message() {
  const [text, setText] = useState("")

  return (
    <>
      <p>{formatDisplayText(text)}</p>
      <input value={text} onChange={(e) => setText(e.target.value)} />
    </>
  )
}

function formatDisplayText(text: string) {
  return text.toUpperCase()
}
```

## Tailwind classes

- Prefer direct inline Tailwind classes in JSX.
- Do not extract class strings into variables or helpers.
- For conditional classes, avoid template literals and ternaries; use `cn`, `clsx`, or `classnames` inline.

```tsx
<input
  className={cn(
    "w-full rounded-lg border px-4 py-2 focus:outline-none focus:ring-4",
    hasWarning && "border-yellow-500 bg-yellow-50 focus:ring-yellow-200",
    hasError && "border-red-500 bg-red-50 focus:ring-red-200"
  )}
/>
```
