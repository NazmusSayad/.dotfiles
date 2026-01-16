---
name: generate-terminal-command
description: This rule explains how to generate a terminal command
---

### Use Linux paths

*Correct:*
```sh
go build -o ./output ./input
go build /c/Users/Sayad/...
```

*Incorrect:*
```sh
go build -o .\\output .\\input
go build C:\\Users\\Sayad\\...
```

### Avoid `cd`

*Correct:*
```sh
ls -l
```

*Incorrect:*
```sh
cd ./path; ls -l
```
