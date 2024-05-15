package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	//import gopherjs
	"github.com/gopherjs/gopherjs/js"
)

const (
	initialTapeSize   = 30000
	tapeExpansionSize = 10000
)

type command byte

const (
	cmdIncPtr command = iota
	cmdDecPtr
	cmdIncByte
	cmdDecByte
	cmdOutput
	cmdInput
	cmdJumpForward
	cmdJumpBackward
)

type instruction struct {
	cmd  command
	arg  int
	jump int
}

func debugPrintln(debugMode bool, a ...interface{}) {
	if debugMode {
		log.Println(a...)
	}
}

func expandTapeIfNeeded(tape []byte, ptr int) []byte {
	if ptr >= len(tape) {
		tape = append(tape, make([]byte, tapeExpansionSize)...)
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

func compileBFCode(code []byte) []instruction {
	bracketStack := make([]int, 0)
	insts := make([]instruction, 0)

	for i := 0; i < len(code); i++ {
		c := code[i]
		switch c {
		case '>':
			if len(insts) > 0 && insts[len(insts)-1].cmd == cmdIncPtr {
				insts[len(insts)-1].arg++
			} else {
				insts = append(insts, instruction{cmd: cmdIncPtr, arg: 1})
			}
		case '<':
			if len(insts) > 0 && insts[len(insts)-1].cmd == cmdDecPtr {
				insts[len(insts)-1].arg++
			} else {
				insts = append(insts, instruction{cmd: cmdDecPtr, arg: 1})
			}
		case '+':
			if len(insts) > 0 && insts[len(insts)-1].cmd == cmdIncByte {
				insts[len(insts)-1].arg++
			} else {
				insts = append(insts, instruction{cmd: cmdIncByte, arg: 1})
			}
		case '-':
			if len(insts) > 0 && insts[len(insts)-1].cmd == cmdDecByte {
				insts[len(insts)-1].arg++
			} else {
				insts = append(insts, instruction{cmd: cmdDecByte, arg: 1})
			}
		case '.':
			insts = append(insts, instruction{cmd: cmdOutput})
		case ',':
			insts = append(insts, instruction{cmd: cmdInput})
		case '[':
			insts = append(insts, instruction{cmd: cmdJumpForward})
			bracketStack = append(bracketStack, len(insts)-1)
		case ']':
			if len(bracketStack) == 0 {
				log.Fatal("Unmatched ']' at position", i)
			}
			openPos := bracketStack[len(bracketStack)-1]
			bracketStack = bracketStack[:len(bracketStack)-1]
			insts[openPos].jump = len(insts)
			insts = append(insts, instruction{cmd: cmdJumpBackward, jump: openPos})
		}
	}
	if len(bracketStack) != 0 {
		log.Fatal("Unmatched '[' at position", bracketStack[len(bracketStack)-1])
	}
	return insts
}

func executeBFCodeOptimized(insts []instruction, tape []byte, stepByStep, debugMode, noInteraction bool, args []byte, output io.Writer) {
	var ptr int
	argIndex := 0
	inputReader := bufio.NewReader(os.Stdin)
	bufferedOutput := bufio.NewWriter(output)

	for ip := 0; ip < len(insts); ip++ {
		inst := insts[ip]
		tape = expandTapeIfNeeded(tape, ptr)

		switch inst.cmd {
		case cmdIncPtr:
			ptr += inst.arg
		case cmdDecPtr:
			ptr -= inst.arg
			if ptr < 0 {
				log.Fatal("Error: Attempt to move pointer left of the starting position")
			}
		case cmdIncByte:
			tape[ptr] += byte(inst.arg)
		case cmdDecByte:
			tape[ptr] -= byte(inst.arg)
		case cmdOutput:
			bufferedOutput.WriteByte(tape[ptr])
		case cmdInput:
			if argIndex < len(args) {
				tape[ptr] = args[argIndex]
				argIndex++
			} else if noInteraction {
				tape[ptr] = 0
			} else {
				tape[ptr] = readInput(inputReader)
			}
		case cmdJumpForward:
			if tape[ptr] == 0 {
				ip = inst.jump
			}
		case cmdJumpBackward:
			if tape[ptr] != 0 {
				ip = inst.jump
			}
		}
		if stepByStep {
			fmt.Printf("Step %d: Executed %d: Tape: ", ip, inst.cmd)
			for i := ptr - 10; i <= ptr+10; i++ {
				if i >= 0 && i < len(tape) {
					fmt.Printf("%v ", tape[i])
				}
			}
			fmt.Println()
		}
	}
	bufferedOutput.Flush()
}

func executeBFCodeInterpreted(code []byte, tape []byte, stepByStep, debugMode, noInteraction bool, args []byte, output io.Writer) {
	var ptr int
	argIndex := 0
	inputReader := bufio.NewReader(os.Stdin)
	bufferedOutput := bufio.NewWriter(output)

	for codeIndex := 0; codeIndex < len(code); codeIndex++ {
		c := code[codeIndex]
		tape = expandTapeIfNeeded(tape, ptr)

		switch c {
		case '>':
			ptr++
			debugPrintln(debugMode, "Pointer moved to", ptr)
			tape = expandTapeIfNeeded(tape, ptr)
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
			bufferedOutput.WriteByte(tape[ptr])
		case ',':
			if argIndex < len(args) {
				tape[ptr] = args[argIndex]
				argIndex++
			} else if noInteraction {
				tape[ptr] = 0
			} else {
				tape[ptr] = readInput(inputReader)
			}
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
	bufferedOutput.Flush()
}

func main() {
	if js.Global != nil {
		js.Global.Set("executeBrainfuck", func(code string, noInteraction bool, input string) string {
			insts := compileBFCode([]byte(code))
			tape := make([]byte, initialTapeSize)
			args := []byte(input)

			output := &bytes.Buffer{}
			executeBFCodeOptimized(insts, tape, false, false, noInteraction, args, output)

			return output.String()
		})
		return
	}

	stepByStep := flag.Bool("s", false, "Enable step-by-step execution")
	debugMode := flag.Bool("d", false, "Enable debug mode")
	filePath := flag.String("f", "", "Path to Brainfuck code file")
	disableOptimizations := flag.Bool("disable-optimizations", false, "Disable optimizations and use the interpreted mode")
	noInteraction := flag.Bool("no-interaction", false, "Disable user input interaction")
	args := flag.String("args", "", "Arguments to pass to the Brainfuck program (space-separated)")
	flag.Parse()

	if *filePath == "" {
		log.Fatal("No file path provided")
	}

	code, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	argBytes := []byte(*args)
	tape := make([]byte, initialTapeSize)
	output := bufio.NewWriter(os.Stdout)

	if *disableOptimizations {
		executeBFCodeInterpreted(code, tape, *stepByStep, *debugMode, *noInteraction, argBytes, output)
	} else {
		insts := compileBFCode(code)
		executeBFCodeOptimized(insts, tape, *stepByStep, *debugMode, *noInteraction, argBytes, output)
	}
	output.Flush()
}
