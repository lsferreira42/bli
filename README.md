```
 ____  _     ___
| __ )| |   |_ _|
|  _ \| |    | |
| |_) | |___ | |
|____/|_____|___|
One of many brainfuck interpreters out there
```


### Why?

Because i'm learning about compilers and interpreters


## Does it work?

Yes, all you need is golang and if you need the web build, gopherjs

## Features

- Interpret Brainfuck code from either a file or standard input
- Step-by-step execution with tape visualization
- Debug mode for enhanced troubleshooting
- Web-based interface using WebAssembly (gopherjs)

## Prerequisites

- Golang
- Make
- Gopherjs

## Installation

### Native Build

To compile the program natively, simply run:

```shell
$ make build
```

### WebAssembly Build

To compile the program to WebAssembly using Emscripten, run:

```shell
$ make web
```

## Usage

### Native

To interpret Brainfuck code from standard input:

```shell
$ echo " >++++++++[-<+++++++++>]<.>>+>-[+]++>++>+++[>[->+++<<+++>]<<]>-----.>->
+++..+++.>-.<<+[>[+>+]>>]<--------------.>>.+++.------.--------.>+.>+." | ./bli
```

To interpret Brainfuck code from a file:

```shell
$ ./bli -f your_file.bf
```

For step-by-step execution:

```shell
$ ./bli -s -f your_file.bf
```

For debug mode:

```shell
$ ./bli -d -f your_file.bf
```

### Web-based

Open `index.html` in your web browser. Paste your Brainfuck code into the textarea and click "Run".

## How It Works

The interpreter starts by reading the Brainfuck code provided either as a command-line argument or from standard input. It then iterates through each character, performing actions based on the Brainfuck instruction set.

The memory is represented by a tape, a 30,000-cell array initialized with zeros. The Brainfuck code manipulates this tape to perform computations, which can then be output as ASCII characters.

### Brainfuck Instruction Set

| Instruction | Utility         |
|-------------|-----------------|
| `.`         | Output the byte at the data pointer as an ASCII encoded character. |
| `,`         | Accept one byte of input and store its value in the byte at the data pointer. |
| `+`         | Increment the byte at the data pointer. |
| `-`         | Decrement the byte at the data pointer. |
| `<`         | Decrement the data pointer. |
| `>`         | Increment the data pointer. |
| `[`         | Jump past the corresponding `]` if the byte at the data pointer is zero. |
| `]`         | Jump back to the corresponding `[` if the byte at the data pointer is nonzero. |


### TODO: talk about the bytecode process