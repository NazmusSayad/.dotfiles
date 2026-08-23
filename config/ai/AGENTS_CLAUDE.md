- Prefer clarity and simplicity over abstraction. Only introduce variables, functions, helpers, or types when logic becomes **very large** and **extremely complex**, or repetition happens **many times**, to keep code simple and easy to follow.
- Use **explicit logic**: avoid `if true: 1; else: 0`, instead use `if true: 1; if false: 0; else: exception` to reduce ambiguity and prevent implicit fallbacks.

- Do not write comments unless instructed.

- Avoid adding test files unless instructed; perform relevant validation only when necessary, combining checks where possible.

- Prefer dedicated tools for working with files instead of using shell scripts.

- Use `behave` skill before planning, making decisions, writing code, or making major changes.
