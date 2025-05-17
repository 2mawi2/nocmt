package main

import (
	"fmt"
	"os"
	"path/filepath"

	"nocmt/processor"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run extract_expected.go <input_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := filepath.Join(filepath.Dir(inputFile), "expected.go")

	content, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	proc := processor.NewGoProcessor(true)
	result, err := proc.StripComments(string(content))
	if err != nil {
		fmt.Printf("Error processing: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(outputFile, []byte(result), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Generated expected output:", outputFile)
}
