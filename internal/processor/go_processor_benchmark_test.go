package processor

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

func BenchmarkGoStripComments(b *testing.B) {
	smallCode := `package main

// This is a comment
func main() {
	// Another comment
	fmt.Println("Hello, World!")  // End of line comment
}
`

	mediumCode := `package main

// Comment 1
/* Block comment
   spanning multiple lines */
import (
	"fmt"  // Import fmt
	"strings"  // Import strings
)

// Main function
func main() {
	// Variable declaration
	greeting := "Hello"  // String variable
	name := "World"      /* Another variable */
	
	// Concatenate strings
	message := fmt.Sprintf("%s, %s!", greeting, name)
	
	/* This is a
	   multiline block comment
	   with multiple lines */
	fmt.Println(message)  // Print the message
	
	// More comments
	// More comments
	// More comments
}
`

	var largeCodeBuilder strings.Builder
	largeCodeBuilder.WriteString("package main\n\n")
	largeCodeBuilder.WriteString("import (\n\t\"fmt\"\n\t\"strings\"\n\t\"math\"\n)\n\n")

	for i := 0; i < 100; i++ {
		largeCodeBuilder.WriteString(fmt.Sprintf("// Function %d documentation\n", i))
		largeCodeBuilder.WriteString(fmt.Sprintf("/* Function %d\n   does something important */\n", i))
		largeCodeBuilder.WriteString(fmt.Sprintf("func function%d() {\n", i))
		largeCodeBuilder.WriteString("\t// Local variable\n")
		largeCodeBuilder.WriteString(fmt.Sprintf("\tval := %d // Value is %d\n", i, i))
		largeCodeBuilder.WriteString("\t/* Block comment inside function */\n")
		largeCodeBuilder.WriteString("\tfmt.Println(val) // Print the value\n")
		largeCodeBuilder.WriteString("}\n\n")
	}
	largeCode := largeCodeBuilder.String()

	benchmarks := []struct {
		name string
		code string
	}{
		{"SmallCode", smallCode},
		{"MediumCode", mediumCode},
		{"LargeCode", largeCode},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			factory := NewProcessorFactory()
			processor, err := factory.GetProcessor("go")
			if err != nil {
				b.Fatal(err)
			}
			for i := 0; i < b.N; i++ {
				_, err := processor.StripComments(bm.code)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkGoFindCommentNodes(b *testing.B) {
	parser := sitter.NewParser()
	parser.SetLanguage(golang.GetLanguage())

	code := `package main

// Line comment
/* Block comment */
func main() {
	// Indented comment
	fmt.Println("Hello")  // End-of-line comment
	/* Multi-line
	   block comment */
}
`

	sourceBytes := []byte(code)
	tree, err := parser.ParseCtx(context.Background(), nil, sourceBytes)
	if err != nil {
		b.Fatal(err)
	}
	defer tree.Close()

	rootNode := tree.RootNode()
	if rootNode == nil {
		b.Fatal("failed to get root node")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		findCommentNodes(rootNode, code)
	}
}

func BenchmarkGoRemoveComments(b *testing.B) {
	code := `package main

// Line comment 1
// Line comment 2
func main() {
	/* Block comment 1 */
	fmt.Println("Hello")  // End-of-line comment
	/* Multi-line
	   block comment 2 */
	// Final comment
}
`

	parser := sitter.NewParser()
	parser.SetLanguage(golang.GetLanguage())

	sourceBytes := []byte(code)
	tree, err := parser.ParseCtx(context.Background(), nil, sourceBytes)
	if err != nil {
		b.Fatal(err)
	}
	defer tree.Close()

	rootNode := tree.RootNode()
	if rootNode == nil {
		b.Fatal("failed to get root node")
	}

	ranges := findCommentNodes(rootNode, code)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		removeComments(code, ranges)
	}
}

func BenchmarkGoRealWorldCode(b *testing.B) {
	var codeBuilder strings.Builder

	codeBuilder.WriteString("// Package processor provides utilities for processing Go source code.\n")
	codeBuilder.WriteString("// It includes functions for removing comments and preserving structure.\n")
	codeBuilder.WriteString("package processor\n\n")

	codeBuilder.WriteString("import (\n")
	codeBuilder.WriteString("\t\"context\" // For context handling\n")
	codeBuilder.WriteString("\t\"fmt\" // For formatting output\n")
	codeBuilder.WriteString("\t\"strings\" // For string manipulation\n\n")

	codeBuilder.WriteString("\t/* External imports */\n")
	codeBuilder.WriteString("\tsitter \"github.com/smacker/go-tree-sitter\" // Tree-sitter parser\n")
	codeBuilder.WriteString("\t\"github.com/smacker/go-tree-sitter/golang\" // Go language support\n")
	codeBuilder.WriteString(")\n\n")

	codeBuilder.WriteString("// ProcessSource processes the given source code by removing comments.\n")
	codeBuilder.WriteString("// It returns the processed code or an error if parsing fails.\n")
	codeBuilder.WriteString("/* This is the main entry point for the package */\n")
	codeBuilder.WriteString("func ProcessSource(source string) (string, error) {\n")
	codeBuilder.WriteString("\t// Initialize parser\n")
	codeBuilder.WriteString("\tparser := sitter.NewParser()\n")
	codeBuilder.WriteString("\tparser.SetLanguage(golang.GetLanguage()) // Set language to Go\n\n")

	codeBuilder.WriteString("\t// Convert source to byte array\n")
	codeBuilder.WriteString("\tsourceBytes := []byte(source) /* Source as bytes */\n\n")

	codeBuilder.WriteString("\t// Parse source code\n")
	codeBuilder.WriteString("\ttree, err := parser.ParseCtx(context.Background(), nil, sourceBytes)\n")
	codeBuilder.WriteString("\tif err != nil {\n")
	codeBuilder.WriteString("\t\treturn \"\", fmt.Errorf(\"failed to parse source code: %w\", err) // Return error\n")
	codeBuilder.WriteString("\t}\n")
	codeBuilder.WriteString("\tif tree == nil { // Check for nil tree\n")
	codeBuilder.WriteString("\t\treturn \"\", fmt.Errorf(\"failed to parse source code\")\n")
	codeBuilder.WriteString("\t}\n")
	codeBuilder.WriteString("\tdefer tree.Close() // Ensure tree is closed\n\n")

	codeBuilder.WriteString("\t/* Get the root node of the AST */\n")
	codeBuilder.WriteString("\trootNode := tree.RootNode()\n")
	codeBuilder.WriteString("\tif rootNode == nil {\n")
	codeBuilder.WriteString("\t\treturn \"\", fmt.Errorf(\"failed to get root node\")\n")
	codeBuilder.WriteString("\t}\n\n")

	codeBuilder.WriteString("\t// Find all comment nodes\n")
	codeBuilder.WriteString("\tcommentRanges := findCommentNodes(rootNode, source)\n\n")

	codeBuilder.WriteString("\t// Remove comments from source\n")
	codeBuilder.WriteString("\tcleanedCode := removeComments(source, commentRanges)\n\n")

	codeBuilder.WriteString("\treturn cleanedCode, nil // Return processed code\n")
	codeBuilder.WriteString("}\n")

	for i := 1; i <= 3; i++ {
		codeBuilder.WriteString(fmt.Sprintf("\n// helperFunction%d is a helper function for processing.\n", i))
		codeBuilder.WriteString(fmt.Sprintf("func helperFunction%d(input string) string {\n", i))
		codeBuilder.WriteString("\t// Local implementation\n")
		codeBuilder.WriteString("\tresult := strings.TrimSpace(input) // Trim spaces\n\n")
		codeBuilder.WriteString("\t/* Process the string */\n")
		codeBuilder.WriteString("\tresult = strings.ReplaceAll(result, \"\\t\", \"  \") // Replace tabs with spaces\n\n")
		codeBuilder.WriteString("\treturn result // Return result\n")
		codeBuilder.WriteString("}\n")
	}

	realWorldCode := codeBuilder.String()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		factory := NewProcessorFactory()
		processor, err := factory.GetProcessor("go")
		if err != nil {
			b.Fatal(err)
		}
		_, err = processor.StripComments(realWorldCode)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoStripCommentsPreserveDirectives(b *testing.B) {
	code := `package main

//go:generate go run gen.go
//go:noinline
// Regular comment
func main() {
	// Another regular comment
	//go:inline
	fmt.Println("Hello")  // End of line comment
	
	/* Block comment with directives inside
	//go:generate protoc --go_out=. *.proto
	*/
}

//go:build linux && (amd64 || arm64)
// +build linux,amd64 linux,arm64
`

	factory := NewProcessorFactory()
	factory.SetPreserveDirectives(true)
	processor, err := factory.GetProcessor("go")
	if err != nil {
		b.Fatalf("Failed to get Go processor: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := processor.StripComments(code)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoComparisonNoDirectives(b *testing.B) {
	code := `package main

// Regular comment
/* Block comment */
func main() {
	// Function comment
	fmt.Println("Hello")  // End of line comment
}
`

	factory := NewProcessorFactory()

	b.Run("StripCommentsNoDirectives", func(b *testing.B) {
		factory.SetPreserveDirectives(false)
		processor, err := factory.GetProcessor("go")
		if err != nil {
			b.Fatalf("Failed to get Go processor: %v", err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := processor.StripComments(code)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("StripCommentsWithDirectives", func(b *testing.B) {
		factory.SetPreserveDirectives(true)
		processor, err := factory.GetProcessor("go")
		if err != nil {
			b.Fatalf("Failed to get Go processor: %v", err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := processor.StripComments(code)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkGoComparisonWithDirectives(b *testing.B) {
	code := `//go:build linux || darwin
// +build linux darwin

package main

//go:generate protoc --go_out=. *.proto
// Regular comment
func main() {
	//go:noinline
	// Function comment
	fmt.Println("Hello")  // End of line comment
}
`
	factory := NewProcessorFactory()

	b.Run("StripCommentsNoDirectives", func(b *testing.B) {
		factory.SetPreserveDirectives(false)
		processor, err := factory.GetProcessor("go")
		if err != nil {
			b.Fatalf("Failed to get Go processor: %v", err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := processor.StripComments(code)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("StripCommentsWithDirectives", func(b *testing.B) {
		factory.SetPreserveDirectives(true)
		processor, err := factory.GetProcessor("go")
		if err != nil {
			b.Fatalf("Failed to get Go processor: %v", err)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := processor.StripComments(code)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkGoParserReuse(b *testing.B) {
	samples := []string{
		`package main

func main() {
	// Comment
	fmt.Println("Hello")
}`,
		`package main

import "fmt"

// Main function
func main() {
	/* Block comment */
	fmt.Println("Hello, World!")
}`,
		`package main

import (
	"fmt"
	"strings"
)

// Function to process string
func process(s string) string {
	// Trim space
	s = strings.TrimSpace(s)
	
	/* Replace special characters
	   with their escaped versions */
	return s
}

func main() {
	// Comment 1
	input := "  test string  "
	// Comment 2
	result := process(input)
	// Comment 3
	fmt.Println(result)
}`,
	}

	b.Run("CreateParserEachTime", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sample := samples[i%len(samples)]

			parser := sitter.NewParser()
			parser.SetLanguage(golang.GetLanguage())

			sourceBytes := []byte(sample)
			tree, err := parser.ParseCtx(context.Background(), nil, sourceBytes)
			if err != nil {
				b.Fatal(err)
			}

			rootNode := tree.RootNode()
			ranges := findCommentNodes(rootNode, sample)
			removeComments(sample, ranges)

			tree.Close()
		}
	})

	b.Run("ReuseParser", func(b *testing.B) {
		parser := sitter.NewParser()
		parser.SetLanguage(golang.GetLanguage())

		for i := 0; i < b.N; i++ {
			sample := samples[i%len(samples)]

			sourceBytes := []byte(sample)
			tree, err := parser.ParseCtx(context.Background(), nil, sourceBytes)
			if err != nil {
				b.Fatal(err)
			}

			rootNode := tree.RootNode()
			ranges := findCommentNodes(rootNode, sample)
			removeComments(sample, ranges)

			tree.Close()
		}
	})
}

func BenchmarkGoMemoryUsage(b *testing.B) {
	code := `package main

import (
	"fmt"
	"strings"
)

// This function does something important
/* It has multiple comment styles
   with varying lengths and formats */
func doSomething() string {
	// Local variable
	result := "processed"  // Assign value
	
	/* Another block comment
	 * with formatted content
	 * spanning multiple lines
	 */
	
	// A group of comments
	// One after another
	// To be processed together
	
	return result // Return the value
}

// Main function for the program
func main() {
	// Call the function
	value := doSomething()
	
	// Print the result
	fmt.Println(value)
}
`

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		factory := NewProcessorFactory()
		processor, err := factory.GetProcessor("go")
		if err != nil {
			b.Fatal(err)
		}
		_, err = processor.StripComments(code)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoFileIO(b *testing.B) {
	b.Skip("File I/O benchmark skipped - enable manually if needed")

	const (
		inputFilePath  = "../testfiles/sample.go"
		outputFilePath = "../testfiles/sample_out.go"
	)

	b.Run("ReadProcessWrite", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			content, err := os.ReadFile(inputFilePath)
			if err != nil {
				b.Fatal(err)
			}

			factory := NewProcessorFactory()
			processor, err := factory.GetProcessor("go")
			if err != nil {
				b.Fatal(err)
			}
			processed, err := processor.StripComments(string(content))
			if err != nil {
				b.Fatal(err)
			}

			err = os.WriteFile(outputFilePath, []byte(processed), 0644)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkGoFullPipeline(b *testing.B) {
	code := `// Copyright 2023 Example Inc.
// All rights reserved.
// License information...

package example

import (
	"fmt"        // For formatted output
	"io"         // I/O interfaces
	"os"         // Operating system interface
	"strings"    // String utilities
)

// Constants used throughout the package
const (
	// Default buffer size
	DefaultBufferSize = 4096 // 4KB buffer
	
	/* Minimum and maximum sizes */
	MinSize = 256
	MaxSize = 1048576 // 1MB
)

// Config stores processing configuration
type Config struct {
	// Input/output settings
	InputFile  string // Source file
	OutputFile string // Destination file
	
	/* Processing options */
	BufferSize int    // Custom buffer size
	Verbose    bool   // Enable verbose logging
}

// NewConfig creates a default configuration
// with reasonable default values
func NewConfig() *Config {
	// Create a new config with defaults
	return &Config{
		BufferSize: DefaultBufferSize,
		Verbose:    false,
	}
}

// Process handles the file processing logic
// It reads from input, processes, and writes to output
func Process(cfg *Config) error {
	// Validate config
	if cfg.InputFile == "" {
		return fmt.Errorf("input file not specified")
	}
	
	// Open input file
	input, err := os.Open(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer input.Close() // Ensure file is closed
	
	// Create output file
	var output io.Writer
	if cfg.OutputFile != "" {
		// If output file specified, open/create it
		outFile, err := os.Create(cfg.OutputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer outFile.Close()
		output = outFile
	} else {
		// Default to stdout
		output = os.Stdout
	}
	
	// Read and process file content
	content, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	
	// Process content (example transformation)
	processedContent := strings.ReplaceAll(string(content), "TODO", "DONE")
	
	// Write processed content
	_, err = fmt.Fprint(output, processedContent)
	return err
}

// Helper function to check if file exists
func fileExists(path string) bool {
	// Check if path exists
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Main function for standalone operation
func main() {
	// Initialize configuration
	cfg := NewConfig()
	
	// Parse command line arguments
	// (implementation omitted for brevity)
	
	// Process files
	err := Process(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}`

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		factory := NewProcessorFactory()
		processor, err := factory.GetProcessor("go")
		if err != nil {
			b.Fatal(err)
		}
		_, err = processor.StripComments(code)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoVeryLargeFiles(b *testing.B) {
	fileSizeConfigurations := []struct {
		testName    string
		targetLines int
	}{
		{"1000Lines", 1000},
		{"2000Lines", 2000},
		{"5000Lines", 5000},
	}

	for _, sizeConfig := range fileSizeConfigurations {
		b.Run(sizeConfig.testName, func(b *testing.B) {
			generatedGoCode := generateRealisticGoCodeWithTargetLineCount(sizeConfig.targetLines)

			goProcessor := createGoProcessorForBenchmark(b)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := goProcessor.StripComments(generatedGoCode)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func generateRealisticGoCodeWithTargetLineCount(targetLineCount int) string {
	var codeBuilder strings.Builder

	addPackageDeclarationAndImports(&codeBuilder)

	functionsToGenerate := calculateFunctionCountForTargetLines(targetLineCount)
	for i := 0; i < functionsToGenerate; i++ {
		addBusinessLogicFunction(&codeBuilder, i)
	}

	return codeBuilder.String()
}

func addPackageDeclarationAndImports(builder *strings.Builder) {
	builder.WriteString("package main\n\n")
	builder.WriteString("import (\n")
	builder.WriteString("\t\"fmt\"\n")
	builder.WriteString("\t\"strings\"\n")
	builder.WriteString("\t\"context\"\n")
	builder.WriteString(")\n\n")
}

func calculateFunctionCountForTargetLines(targetLines int) int {
	const averageLinesPerFunction = 10
	return targetLines / averageLinesPerFunction
}

func addBusinessLogicFunction(builder *strings.Builder, functionIndex int) {
	fmt.Fprintf(builder, "// Function %d performs important business logic\n", functionIndex)
	fmt.Fprintf(builder, "// This function handles case %d specifically\n", functionIndex)
	fmt.Fprintf(builder, "func processData%d(input string) (string, error) {\n", functionIndex)
	builder.WriteString("\t// Validate input parameters\n")
	builder.WriteString("\tif input == \"\" {\n")
	builder.WriteString("\t\t// Return error for empty input\n")
	builder.WriteString("\t\treturn \"\", fmt.Errorf(\"empty input\")\n")
	builder.WriteString("\t}\n")
	builder.WriteString("\t\n")
	builder.WriteString("\t// Process the input data\n")
	fmt.Fprintf(builder, "\tresult := fmt.Sprintf(\"processed_%%s_%d\", input) // Format result\n", functionIndex)
	builder.WriteString("\treturn result, nil // Return success\n")
	builder.WriteString("}\n\n")
}

func createGoProcessorForBenchmark(b *testing.B) LanguageProcessor {
	factory := NewProcessorFactory()
	processor, err := factory.GetProcessor("go")
	if err != nil {
		b.Fatal(err)
	}
	return processor
}
