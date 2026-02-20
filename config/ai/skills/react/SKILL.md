---
name: react
description: Applies common React style guides and conventions. Use when writing or reviewing React components, hooks, or JSX, or when the user asks for React style or best practices.
---

## 1. Don't over-declare

Do not over-declare variables, functions, or components unless the logic is extremely complex. Prefer inlining when it's simple.

❌ Incorrect:

```tsx
import { useState } from "react";

function MyComponent() {
  const [count, setCount] = useState(0);

  function handleIncrement() {
    setCount(count + 1);
  }

  return (
    <div>
      <p>Count: {count}</p>
      <button onClick={handleIncrement}>Increment</button>
    </div>
  );
}
```

Inline it when it's not extremely complex.

✅ Correct:

```tsx
import { useState } from "react";

function MyComponent() {
  const [count, setCount] = useState(0);

  return (
    <div>
      <p>Count: {count}</p>
      <button onClick={() => setCount(count + 1)}>Increment</button>
    </div>
  );
}
```

You may declare functions when the logic is extremely complex.

✅ Correct:

```tsx
import { useState } from "react";

function MyComponent() {
  const [count, setCount] = useState(0);

  function handleIncrement() {
    // YOUR EXTREMELY COMPLEX LOGIC HERE
    setCount(count + 1);
  }

  return (
    <div>
      <p>Count: {count}</p>
      <button onClick={handleIncrement}>Increment</button>
    </div>
  );
}
```

## 2. Keep pure functions outside the component

Pure utility functions must be declared outside React components.

❌ Incorrect:

```tsx
import { useState } from "react";

function MyComponent() {
  const [text, setText] = useState("");

  function processDisplayText(text: string) {
    // ...
    return text.toUpperCase();
  }

  return (
    <div>
      <p>Text: {processDisplayText(text)}</p>
      <input onChange={(e) => setText(e.target.value)} />
    </div>
  );
}
```

✅ Correct:

```tsx
import { useState } from "react";

function processDisplayText(text: string) {
  // ...
  return text.toUpperCase();
}

function MyComponent() {
  const [text, setText] = useState("");

  return (
    <div>
      <p>Text: {processDisplayText(text)}</p>
      <input onChange={(e) => setText(e.target.value)} />
    </div>
  );
}
```
