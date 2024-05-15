package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	initialTapeSize   = 30000
	tapeExpansionSize = 10000
)

func debugPrintln(debugMode bool, a ...interface{}) {
	if debugMode {
		log.Println(a...)
	}
}

func expandTapeIfNeeded(tape []byte, ptr int, debugMode bool) []byte {
	if ptr >= len(tape) {
		tape = append(tape, make([]byte, tapeExpansionSize)...)
		debugPrintln(debugMode, "Tape expanded to", len(tape))
	}
	return tape
}

func readInput(inputReader *bufio.Reader) byte {
	char, _, err := inputReader.ReadRune()
	if err != nil {
		if err == io.EOF {
			return 0
		}
		log.Fatal("Input read error:", err)
	}
	return byte(char)
}

func executeBFCode(code []byte, tape []byte, stepByStep, debugMode bool) {
	var ptr int
	output := make([]byte, 0, initialTapeSize)
	inputReader := bufio.NewReader(os.Stdin)

	for codeIndex := 0; codeIndex < len(code); codeIndex++ {
		c := code[codeIndex]
		tape = expandTapeIfNeeded(tape, ptr, debugMode)

		switch c {
		case '>':
			ptr++
			debugPrintln(debugMode, "Pointer moved to", ptr)
			tape = expandTapeIfNeeded(tape, ptr, debugMode)
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
			tape[ptr] = readInput(inputReader)
		case '[':
			if tape[ptr] == 0 {
				balance := 1
				for balance > 0 {
					codeIndex++
					if codeIndex >= len(code) {
						log.Fatal("Unmatched '[' bracket")
					}
					switch code[codeIndex] {
					case '[':
						balance++
					case ']':
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
					switch code[codeIndex] {
					case ']':
						balance++
					case '[':
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
	stepByStep := flag.Bool("s", false, "Enable step-by-step execution")
	debugMode := flag.Bool("d", false, "Enable debug mode")
	filePath := flag.String("f", "", "Path to Brainfuck code file")
	flag.Parse()

	if *filePath == "" {
		log.Fatal("No file path provided")
	}

	code, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	tape := make([]byte, initialTapeSize)
	executeBFCode(code, tape, *stepByStep, *debugMode)
}
