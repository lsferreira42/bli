# Simple and lazy Brainf*ck interpreter


### Why?

Because i'm learning about compilers and interpreters


## Does it work?

Yes, all you need is GCC and make ( in fact you just need gcc, but typing make is faster )

```shell
$ make
gcc -Wall -Werror bli.c -o bli
$ ./bli "++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++."
Hello World!
```

And that's it!!

### Why Brainf*ck?

As much as Brainfuck can be a little bit intimidating at first look, it's a very simple esoteric language , it has a VERY compact instruction set ( 8 to be exact ) and it's functioning is kinda funny.


### How it works?

it starting by reading the first parameter passed throug the command line, than it loop through it char by char and do a action on the memory tape.

Also, the tape is a 30000 array with 0 on every position, the instruction set move ou update the value of a position on this tape that's later converted to it's equivalent char in the ASCII table, talking about tables, that's one containing all possible brainfuck instructions:

| Instruction | Utility |
|------|-------|
| "." (dot) | print the current cell value    |
| "+" and "-" | increment or decrement the value of the current cell |
| "," | read a char from the user and set at the pointer position |
| "[" and "]" | start and end a loop when the current cell is not zero |
| "<" and ">" | move the point to the right or left |

### Visual explaining: TODO

Add a visual run through the hello world example