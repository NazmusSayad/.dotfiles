---
name: react
description: This skill guides how to write good React code.
---

# React Styles

## Do not over declare variables, functions, or components unless it's extremely complex.

❌ Incorrect:

```jsx
import { useState } from 'react'

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

#### Just inline it when it's not extremely complex.

✅ Correct:

```jsx
import { useState } from 'react'

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

#### Can declare functions when it's extremely complex.

✅ Correct:

```jsx
import { useState } from 'react'

function MyComponent() {
  const [count, setCount] = useState(0)

  function handleIncrement() {
    // YOUR EXTREMELY COMPLEX LOGIC HERE
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

## Put pure functions always outside of the component.

❌ Incorrect:

```jsx
import { useState } from 'react'

function MyComponent() {
  const [text, setText] = useState('')

  function processDisplayText(text: string) {
    // ...
    return text.toUpperCase()
  }

  return (
    <div>
      <p>Text: {processDisplayText(text)}</p>
      <input onChange={(e) => setText(e.target.value)} />
    </div>
  )
}
```

✅ Correct:

```jsx
import { useState } from 'react'

function processDisplayText(text: string) {
  // ...
  return text.toUpperCase()
}

function MyComponent() {
  const [text, setText] = useState('')

  return (
    <div>
      <p>Text: {processDisplayText(text)}</p>
      <input onChange={(e) => setText(e.target.value)} />
    </div>
  )
}
```
