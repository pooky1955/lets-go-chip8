# lets-go-chip8
Implementing the Chip8 Interpreter in Golang, rendered in the terminal

Here's a demo below:
[![asciicast](https://asciinema.org/a/x7FrZRHnGfmG0hgYqtPKYacdn.svg)](https://asciinema.org/a/x7FrZRHnGfmG0hgYqtPKYacdn)

Currently, there is no keyboard or speaker support, but it passes the following tests:
- Chip-8 Splash Screen
- IBM Logo
- Corax+ opcode test
- Flags test

The tests are taken from Timendus' Test Suite, available at https://github.com/Timendus/chip8-test-suite


Difficulties encountered
---

During the development of the Chip8 Interpreter, I was writing to the VF register before performing operations such as addition or subtraction, while the expected behaviour should set the VF register after the operation. This is to handle the case where you use the VF register as one of your operands.

The other two bugs I made during the development was confusing row,column with x,y and forgetting to read the Opcode with binary.BigEndian.Uint16. As well, when reading the draw opcode DXYN, I didn't extract X and Y properly - I mistakenly read the first nibble to X, and the second nibble to Y instead of the second and third respectively.

Installation
---

There are no dependencies, you can run it simply with
```bash
go build .
./chip8 1-chip8-logo.ch8
