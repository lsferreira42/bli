package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func stringToBrainfuck(input string) string {
	var result strings.Builder
	var previousChar byte

	for i := 0; i < len(input); i++ {
		currentChar := input[i]
		diff := int(currentChar - previousChar)
		previousChar = currentChar

		if diff > 0 {
			result.WriteString(strings.Repeat("+", diff))
		} else if diff < 0 {
			result.WriteString(strings.Repeat("-", -diff))
		}
		result.WriteString(".")
	}

	return result.String()
}

func readStdin() ([]byte, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeCharDevice != 0 {
		return nil, nil
	}
	reader := bufio.NewReader(os.Stdin)
	return io.ReadAll(reader)
}

func main() {
	input, err := readStdin()
	if err != nil {
		fmt.Println(err)
		return
	}
	brainfuckCode := stringToBrainfuck(string(input))
	fmt.Println(brainfuckCode)
}
