# lets-go-chip8
Implementing the Chip8 Interpreter in Golang, rendered in the terminal

Here's a demo below:
[![asciicast](https://asciinema.org/a/x7FrZRHnGfmG0hgYqtPKYacdn.svg)](https://asciinema.org/a/x7FrZRHnGfmG0hgYqtPKYacdn)

Currently, there is no keyboard or speaker support, but it passes the following tests:
- Chip-8 Splash Screen
- IBM Logo
- Corax+ opcode test
- Flags test

Installation
---

There are no dependencies, you can run it simply with
```bash
go build .
./chip8 1-chip8-logo.ch8
