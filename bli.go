package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func printDebug(format string, a ...interface{}) {
	if _, exists := os.LookupEnv("DEBUG"); exists {
		log.Printf(format, a...)
	}
}

func resizeTape(tape *[]byte, ptr *int) {
	oldSize := len(*tape)
	newSize := oldSize * 2
	newTape := make([]byte, newSize)
	copy(newTape, *tape)
	*tape = newTape
	*ptr += oldSize
	printDebug("Resized tape to %d bytes\n", newSize)
}

func executeBFCode(code []byte, tape []byte, stepByStep, debugMode bool) {
	var ptr int
	output := make([]byte, 300000)
	var outputIndex int

	for codeIndex := 0; codeIndex < len(code); codeIndex++ {
		c := code[codeIndex]

		if ptr >= len(tape) {
			resizeTape(&tape, &ptr)
		}

		switch c {
		case '>':
			ptr++
		case '<':
			ptr--
		case '+':
			tape[ptr]++
		case '-':
			tape[ptr]--
		case '.':
			output[outputIndex] = tape[ptr]
			outputIndex++
		case ',':
			input := bufio.NewReader(os.Stdin)
			char, _, err := input.ReadRune()
			if err != nil {
				log.Fatal("Failed to read input\n")
			}
			tape[ptr] = byte(char)
		case '[':
			if tape[ptr] == 0 {
				balance := 1
				for balance != 0 {
					codeIndex++
					if codeIndex >= len(code) {
						log.Fatal("Jumped to unbalanced bracket\n")
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
				for balance != 0 {
					codeIndex--
					if codeIndex < 0 {
						log.Fatal("Jumped to unbalanced bracket\n")
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
	}

	output = output[:outputIndex] // Truncate to actual size
	fmt.Println(string(output))
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
			os.Setenv("DEBUG", "true")
		case "-c":
			if i+1 < len(args) {
				filePath = args[i+1]
				i++
			} else {
				log.Fatal("Expected file path after -c\n")
			}
		}
	}

	if debugMode {
		printDebug("File path: %s\n", filePath)
	}

	var code []byte
	if filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Error reading file: %v\n", err)
		}
		code = data
	} else {
		reader := bufio.NewReader(os.Stdin)
		var err error
		code, err = io.ReadAll(reader)
		if err != nil {
			log.Fatalf("Error reading from stdin: %v\n", err)
			return
		}
	}

	tapeSize := 30000
	tape := make([]byte, tapeSize)
	executeBFCode(code, tape, stepByStep, debugMode)
}
