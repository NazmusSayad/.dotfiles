## Rules

- NEVER run `go build` or `go run`. Ask the user to run those commands themselves.
- NEVER run repository sync/update scripts directly (.exe). Ask the user to run them.
- If the user explicitly said to run then only run via `go run <script>`
