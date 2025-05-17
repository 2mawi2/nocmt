package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	lines := []string{
		"//go:build linux && !windows",
		"// +build linux,!windows",
		"",
		"package main",
		"import (",
		"\t\"fmt\"",
		")",
		"",
		"",
		"const Version = \"v1.0.0\"",
		"",
		"",
		"var (",
		"",
		"\tname =  \"Gopher\"",
		"\tage = 10",
		")",
		"",
		"",
		"",
		"",
		"func hello() {",
		"\tfmt.Println(\"Hello\")",
		"}",
		"",
		"func main() {",
		"\thello()",
		"",
		"\t//go:generate echo \"generate something\"",
		"\t//go:noinline",
		"\tif true {",
		"\t\tfmt.Println(\"Conditional\")",
		"\t} else  {",
		"\t\tfmt.Println(\"Else branch\")",
		"\t}",
		"",
		"",
		"",
		"}",
	}

	result := strings.Join(lines, "\n") + "\n"
	err := os.WriteFile("expected.go", []byte(result), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Generated expected.go")
}
