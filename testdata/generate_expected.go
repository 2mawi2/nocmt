// generate_expected.go - Utility to generate expected output from test fixtures
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"nocmt/processor"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run generate_expected.go <language>")
		fmt.Println("       go run generate_expected.go all")
		fmt.Println("\nAvailable languages: go, javascript, typescript, python, rust, swift, kotlin, css, csharp, bash, java")
		os.Exit(1)
	}

	language := os.Args[1]

	if language == "all" {
		languages := []string{"go", "javascript", "typescript", "python", "rust", "swift", "kotlin", "css", "csharp", "bash", "java"}
		for _, lang := range languages {
			if err := processLanguage(lang); err != nil {
				fmt.Printf("Error processing %s: %v\n", lang, err)
			}
		}
		return
	}

	if err := processLanguage(language); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func processLanguage(language string) error {
	var proc processor.LanguageProcessor
	var originalPath, expectedPath, ext string

	switch strings.ToLower(language) {
	case "go":
		proc = processor.NewGoProcessor(true)
		ext = "go"
	case "javascript":
		proc = processor.NewJavaScriptProcessor(true)
		ext = "js"
	case "typescript":
		proc = processor.NewTypeScriptProcessor(true)
		ext = "ts"
	case "python":
		proc = processor.NewPythonProcessor(true)
		ext = "py"
	case "rust":
		proc = processor.NewRustProcessor(true)
		ext = "rs"
	case "swift":
		proc = processor.NewSwiftProcessor(true)
		ext = "swift"
	case "kotlin":
		proc = processor.NewKotlinProcessor(true)
		ext = "kt"
	case "css":
		proc = processor.NewCSSProcessor(true)
		ext = "css"
	case "csharp":
		proc = processor.NewCSharpProcessor(true)
		ext = "cs"
	case "bash":
		proc = processor.NewBashProcessor(true)
		ext = "sh"
	case "java":
		proc = processor.NewJavaProcessor(true)
		ext = "java"
	default:
		return fmt.Errorf("unknown language: %s", language)
	}

	originalPath = filepath.Join("testdata", language, fmt.Sprintf("original.%s", ext))
	expectedPath = filepath.Join("testdata", language, fmt.Sprintf("expected.%s", ext))

	// Read original content
	content, err := os.ReadFile(originalPath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Process content
	processed, err := proc.StripComments(string(content))
	if err != nil {
		return fmt.Errorf("error processing file: %v", err)
	}

	// Write processed content
	err = os.WriteFile(expectedPath, []byte(processed), 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	fmt.Printf("Successfully generated %s\n", expectedPath)
	return nil
}
