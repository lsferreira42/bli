package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func debugPrintln(debugMode bool, a ...interface{}) {
	if debugMode {
		fmt.Println(a...)
	}
}

// TODO: implement the debug mode
func executeBFCode(code []byte, tape []byte, stepByStep, debugMode bool) {
	var ptr int
	output := make([]byte, 0, 30000) // Starts with a capacity of 30000, can grow dynamically
	inputReader := bufio.NewReader(os.Stdin)

	for codeIndex := 0; codeIndex < len(code); codeIndex++ {
		c := code[codeIndex]
		if ptr >= len(tape) {
			tape = append(tape, make([]byte, 10000)...)
		}

		switch c {
		case '>':
			ptr++
			debugPrintln(debugMode, "Pointer moved to", ptr)
			if ptr >= len(tape) {
				tape = append(tape, make([]byte, 10000)...)
				debugPrintln(debugMode, "Tape expanded to", len(tape))
			}
		case '<':
			if ptr > 0 {
				ptr--
			} else {
				log.Fatal("Error: Attempt to move pointer left of the starting position")
			}
		case '+':
			tape[ptr]++
		case '-':
			tape[ptr]--
		case '.':
			output = append(output, tape[ptr])
		case ',':
			char, _, err := inputReader.ReadRune()
			if err != nil {
				if err == io.EOF {
					tape[ptr] = 0
				} else {
					log.Fatal("Input read error:", err)
				}
			} else {
				tape[ptr] = byte(char)
			}
		case '[':
			if tape[ptr] == 0 {
				balance := 1
				for balance > 0 {
					codeIndex++
					if codeIndex >= len(code) {
						log.Fatal("Unmatched '[' bracket")
					}
					if code[codeIndex] == '[' {
						balance++
					} else if code[codeIndex] == ']' {
						balance--
					}
				}
			}
		case ']':
			if tape[ptr] != 0 {
				balance := 1
				for balance > 0 {
					codeIndex--
					if codeIndex < 0 {
						log.Fatal("Unmatched ']' bracket")
					}
					if code[codeIndex] == ']' {
						balance++
					} else if code[codeIndex] == '[' {
						balance--
					}
				}
			}
		}
		if stepByStep {
			fmt.Printf("Step %d: Executed %c: Tape: ", codeIndex, c)
			for i := ptr - 10; i <= ptr+10; i++ {
				if i >= 0 && i < len(tape) {
					fmt.Printf("%v ", tape[i])
				}
			}
			fmt.Println()
		}
	}

	fmt.Print(string(output))
}

func main() {
	var stepByStep, debugMode bool
	var filePath string

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-s":
			stepByStep = true
		case "-d":
			debugMode = true
		case "-f":
			if i+1 < len(args) {
				filePath = args[i+1]
				i++
			} else {
				log.Fatal("Expected filename after -f")
			}
		default:
			log.Fatalf("Unknown option %s\n", args[i])
		}
	}

	// Read code from file or stdin
	var code []byte
	if filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		code = data
	} else {
		log.Fatal("No file path provided")
	}

	// Initialize tape with 30,000 bytes as commonly expected by Brainfuck programs
	tape := make([]byte, 30000)
	executeBFCode(code, tape, stepByStep, debugMode)
}
