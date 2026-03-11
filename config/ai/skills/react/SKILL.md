---
name: react
description: React code style, conventions, best practices, and patterns for building clean, efficient, and scalable applications. Use when writing or reviewing React components, hooks, and JSX.
---

### Don't over-declare

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

### Keep pure functions outside the component

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

### Tailwind Classes

Prefer direct inline Tailwind classes in JSX. Do not extract them into variables or helpers. Only extract when the same classes are reused extensively or the logic becomes very complex.

❌ Incorrect:

```tsx
function MyComponent() {
  const containerClasses = "flex flex-col gap-4 p-6 bg-white rounded-lg";
  const buttonClasses =
    "px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600";

  return (
    <div className={containerClasses}>
      <h1>Title</h1>
      <button className={buttonClasses}>Click me</button>
    </div>
  );
}
```

✅ Correct:

```tsx
function MyComponent() {
  return (
    <div className="flex flex-col gap-4 p-6 bg-white rounded-lg">
      <h1>Title</h1>
      <button className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">
        Click me
      </button>
    </div>
  );
}
```

✅ Correct (complex classes, still direct):

```tsx
function Form({ isSubmitting, hasError }: FormProps) {
  return (
    <form className="space-y-6 rounded-xl border border-gray-200 bg-gradient-to-br from-white to-gray-50 p-8 shadow-sm transition-all duration-200 hover:shadow-md">
      <fieldset
        disabled={isSubmitting}
        className="disabled:opacity-50 disabled:pointer-events-none transition-opacity duration-200"
      >
        <input
          className={`w-full rounded-lg border-2 px-4 py-2 transition-colors ${
            hasError
              ? "border-red-500 bg-red-50 focus:border-red-600 focus:ring-red-200"
              : "border-gray-300 focus:border-blue-500 focus:ring-blue-200"
          } focus:outline-none focus:ring-4`}
        />
      </fieldset>
    </form>
  );
}
```

✅ Correct (using `cn` or similar utility, if available in the project):

```tsx
import { cn } from "@/lib/utils"; // or from "clsx", "classnames"

function Form() {
  return (
    <form className="space-y-6 rounded-xl border border-gray-200 bg-gradient-to-br from-white to-gray-50 p-8 shadow-sm transition-all duration-200 hover:shadow-md">
      <fieldset className="disabled:opacity-50 disabled:pointer-events-none transition-opacity duration-200">
        <input
          className={cn(
            "w-full rounded-lg border-2 px-4 py-2 transition-colors focus:outline-none focus:ring-4 border-gray-300 focus:border-blue-500 focus:ring-blue-200",
            hasWarning &&
              "border-yellow-500 bg-yellow-50 focus:border-yellow-600 focus:ring-yellow-200",
            hasError &&
              "border-red-500 bg-red-50 focus:border-red-600 focus:ring-red-200"
          )}
        />
      </fieldset>
    </form>
  );
}
```
