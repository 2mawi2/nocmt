package processor

import (
	"fmt"
	"strings"
	"testing"
)

type LanguageBenchmarkConfig struct {
	languageName    string
	fileExtension   string
	codeGenerator   func(int) string
	commentPatterns []string
}

func BenchmarkAllLanguagesComparison(b *testing.B) {
	languageConfigs := []LanguageBenchmarkConfig{
		{
			languageName:    "javascript",
			fileExtension:   ".js",
			codeGenerator:   generateJavaScriptCode,
			commentPatterns: []string{"//", "/*", "*/"},
		},
		{
			languageName:    "typescript",
			fileExtension:   ".ts",
			codeGenerator:   generateTypeScriptCode,
			commentPatterns: []string{"//", "/*", "*/"},
		},
		{
			languageName:    "python",
			fileExtension:   ".py",
			codeGenerator:   generatePythonCode,
			commentPatterns: []string{"#"},
		},
		{
			languageName:    "rust",
			fileExtension:   ".rs",
			codeGenerator:   generateRustCode,
			commentPatterns: []string{"//", "/*", "*/"},
		},
		{
			languageName:    "java",
			fileExtension:   ".java",
			codeGenerator:   generateJavaCode,
			commentPatterns: []string{"//", "/*", "*/"},
		},
		{
			languageName:    "kotlin",
			fileExtension:   ".kt",
			codeGenerator:   generateKotlinCode,
			commentPatterns: []string{"//", "/*", "*/"},
		},
		{
			languageName:    "swift",
			fileExtension:   ".swift",
			codeGenerator:   generateSwiftCode,
			commentPatterns: []string{"//", "/*", "*/"},
		},
		{
			languageName:    "cpp",
			fileExtension:   ".cpp",
			codeGenerator:   generateCppCode,
			commentPatterns: []string{"//", "/*", "*/"},
		},
		{
			languageName:    "csharp",
			fileExtension:   ".cs",
			codeGenerator:   generateCSharpCode,
			commentPatterns: []string{"//", "/*", "*/"},
		},
		{
			languageName:    "bash",
			fileExtension:   ".sh",
			codeGenerator:   generateBashCode,
			commentPatterns: []string{"#"},
		},
		{
			languageName:    "css",
			fileExtension:   ".css",
			codeGenerator:   generateCSSCode,
			commentPatterns: []string{"/*", "*/"},
		},
	}

	testSizes := []struct {
		name      string
		lineCount int
	}{
		{"Small_100Lines", 100},
		{"Medium_500Lines", 500},
		{"Large_1000Lines", 1000},
	}

	for _, langConfig := range languageConfigs {
		for _, size := range testSizes {
			benchmarkName := fmt.Sprintf("%s_%s", langConfig.languageName, size.name)

			b.Run(benchmarkName, func(b *testing.B) {
				sourceCode := langConfig.codeGenerator(size.lineCount)

				factory := NewProcessorFactory()
				processor, err := factory.GetProcessor(langConfig.languageName)
				if err != nil {
					b.Fatalf("Failed to get processor for %s: %v", langConfig.languageName, err)
				}

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_, err := processor.StripComments(sourceCode)
					if err != nil {
						b.Fatalf("Failed to strip comments for %s: %v", langConfig.languageName, err)
					}
				}
			})
		}
	}
}

func generateJavaScriptCode(targetLines int) string {
	var builder strings.Builder

	builder.WriteString("// JavaScript module for data processing\n")
	builder.WriteString("/* Multi-line comment\n   describing the module purpose */\n")
	builder.WriteString("const express = require('express');\n")
	builder.WriteString("const { validateInput } = require('./utils');\n\n")

	functionsToGenerate := targetLines / 8
	for i := 0; i < functionsToGenerate; i++ {
		fmt.Fprintf(&builder, "// Function %d: processes user data\n", i)
		fmt.Fprintf(&builder, "function processUserData%d(userData) {\n", i)
		builder.WriteString("  // Validate input parameters\n")
		builder.WriteString("  if (!userData || typeof userData !== 'object') {\n")
		builder.WriteString("    throw new Error('Invalid user data'); // Error handling\n")
		builder.WriteString("  }\n")
		fmt.Fprintf(&builder, "  return { ...userData, processed: true, id: %d }; // Return processed data\n", i)
		builder.WriteString("}\n\n")
	}

	return builder.String()
}

func generateTypeScriptCode(targetLines int) string {
	var builder strings.Builder

	builder.WriteString("// TypeScript interface definitions\n")
	builder.WriteString("/* Application-wide type definitions\n   for better type safety */\n")
	builder.WriteString("interface UserData {\n")
	builder.WriteString("  id: number; // Unique identifier\n")
	builder.WriteString("  name: string; // User display name\n")
	builder.WriteString("  email: string; // Contact email\n")
	builder.WriteString("}\n\n")

	functionsToGenerate := targetLines / 10
	for i := 0; i < functionsToGenerate; i++ {
		fmt.Fprintf(&builder, "// Service class %d for data management\n", i)
		fmt.Fprintf(&builder, "class DataService%d {\n", i)
		builder.WriteString("  // Private data store\n")
		builder.WriteString("  private data: UserData[] = [];\n\n")
		builder.WriteString("  // Method to add user data\n")
		builder.WriteString("  public addUser(user: UserData): void {\n")
		builder.WriteString("    // Validation logic\n")
		builder.WriteString("    this.data.push(user); // Add to collection\n")
		builder.WriteString("  }\n")
		builder.WriteString("}\n\n")
	}

	return builder.String()
}

func generatePythonCode(targetLines int) string {
	var builder strings.Builder

	builder.WriteString("# Python data processing module\n")
	builder.WriteString("# Handles user data validation and transformation\n")
	builder.WriteString("import json\n")
	builder.WriteString("from typing import Dict, List, Optional\n\n")

	functionsToGenerate := targetLines / 8
	for i := 0; i < functionsToGenerate; i++ {
		fmt.Fprintf(&builder, "def process_data_%d(data: Dict) -> Dict:\n", i)
		builder.WriteString("    # Validate input data structure\n")
		builder.WriteString("    if not isinstance(data, dict):\n")
		builder.WriteString("        raise ValueError('Data must be a dictionary')  # Type validation\n")
		builder.WriteString("    \n")
		builder.WriteString("    # Transform and return processed data\n")
		fmt.Fprintf(&builder, "    return {'processed': True, 'original': data, 'processor_id': %d}  # Return result\n", i)
		builder.WriteString("\n")
	}

	return builder.String()
}

func generateRustCode(targetLines int) string {
	var builder strings.Builder

	builder.WriteString("// Rust data processing module\n")
	builder.WriteString("/* Safe and efficient data handling\n   with zero-cost abstractions */\n")
	builder.WriteString("use std::collections::HashMap;\n")
	builder.WriteString("use serde::{Deserialize, Serialize};\n\n")

	functionsToGenerate := targetLines / 12
	for i := 0; i < functionsToGenerate; i++ {
		fmt.Fprintf(&builder, "// Function %d: processes user data safely\n", i)
		fmt.Fprintf(&builder, "pub fn process_user_data_%d(data: &str) -> Result<String, Box<dyn std::error::Error>> {\n", i)
		builder.WriteString("    // Parse JSON input\n")
		builder.WriteString("    let parsed: HashMap<String, serde_json::Value> = serde_json::from_str(data)?;\n")
		builder.WriteString("    \n")
		builder.WriteString("    // Validate required fields\n")
		builder.WriteString("    if !parsed.contains_key(\"id\") {\n")
		builder.WriteString("        return Err(\"Missing required field: id\".into()); // Error handling\n")
		builder.WriteString("    }\n")
		builder.WriteString("    \n")
		builder.WriteString("    // Return processed result\n")
		fmt.Fprintf(&builder, "    Ok(format!(\"{{\\\"processed\\\": true, \\\"processor\\\": %d}}\")) // Success response\n", i)
		builder.WriteString("}\n\n")
	}

	return builder.String()
}

func generateJavaCode(targetLines int) string {
	var builder strings.Builder

	builder.WriteString("// Java data processing service\n")
	builder.WriteString("/* Enterprise-grade data handling\n   with comprehensive error management */\n")
	builder.WriteString("package com.example.processor;\n\n")
	builder.WriteString("import java.util.*;\n")
	builder.WriteString("import java.util.concurrent.ConcurrentHashMap;\n\n")

	functionsToGenerate := targetLines / 15
	for i := 0; i < functionsToGenerate; i++ {
		fmt.Fprintf(&builder, "// Service class %d for data processing\n", i)
		fmt.Fprintf(&builder, "public class DataProcessor%d {\n", i)
		builder.WriteString("    // Thread-safe data storage\n")
		builder.WriteString("    private final Map<String, Object> dataStore = new ConcurrentHashMap<>();\n\n")
		builder.WriteString("    // Method to process user input\n")
		builder.WriteString("    public Map<String, Object> processData(Map<String, Object> input) {\n")
		builder.WriteString("        // Input validation\n")
		builder.WriteString("        if (input == null || input.isEmpty()) {\n")
		builder.WriteString("            throw new IllegalArgumentException(\"Input cannot be null or empty\"); // Validation error\n")
		builder.WriteString("        }\n")
		builder.WriteString("        \n")
		builder.WriteString("        // Process and return result\n")
		fmt.Fprintf(&builder, "        Map<String, Object> result = new HashMap<>(input); // Copy input\n")
		fmt.Fprintf(&builder, "        result.put(\"processed\", true); // Mark as processed\n")
		fmt.Fprintf(&builder, "        result.put(\"processorId\", %d); // Add processor ID\n", i)
		builder.WriteString("        return result;\n")
		builder.WriteString("    }\n")
		builder.WriteString("}\n\n")
	}

	return builder.String()
}

func generateKotlinCode(targetLines int) string {
	var builder strings.Builder

	builder.WriteString("// Kotlin data processing utilities\n")
	builder.WriteString("/* Modern JVM language with concise syntax\n   and null safety features */\n")
	builder.WriteString("package com.example.processor\n\n")
	builder.WriteString("import kotlinx.coroutines.*\n")
	builder.WriteString("import kotlinx.serialization.Serializable\n\n")

	functionsToGenerate := targetLines / 12
	for i := 0; i < functionsToGenerate; i++ {
		fmt.Fprintf(&builder, "// Data class %d for type-safe processing\n", i)
		fmt.Fprintf(&builder, "@Serializable\n")
		fmt.Fprintf(&builder, "data class ProcessedData%d(\n", i)
		builder.WriteString("    val id: String, // Unique identifier\n")
		builder.WriteString("    val data: Map<String, Any>, // Original data\n")
		builder.WriteString("    val processed: Boolean = true // Processing status\n")
		builder.WriteString(")\n\n")
		fmt.Fprintf(&builder, "// Suspend function %d for async processing\n", i)
		fmt.Fprintf(&builder, "suspend fun processDataAsync%d(input: Map<String, Any>): ProcessedData%d {\n", i, i)
		builder.WriteString("    // Simulate async processing\n")
		builder.WriteString("    delay(1) // Non-blocking delay\n")
		builder.WriteString("    \n")
		builder.WriteString("    // Return processed result\n")
		fmt.Fprintf(&builder, "    return ProcessedData%d(\n", i)
		fmt.Fprintf(&builder, "        id = \"proc_%d\", // Generated ID\n", i)
		builder.WriteString("        data = input // Original data\n")
		builder.WriteString("    )\n")
		builder.WriteString("}\n\n")
	}

	return builder.String()
}

func generateSwiftCode(targetLines int) string {
	var builder strings.Builder

	builder.WriteString("// Swift data processing framework\n")
	builder.WriteString("/* iOS/macOS compatible data handling\n   with protocol-oriented programming */\n")
	builder.WriteString("import Foundation\n")
	builder.WriteString("import Combine\n\n")

	functionsToGenerate := targetLines / 14
	for i := 0; i < functionsToGenerate; i++ {
		fmt.Fprintf(&builder, "// Protocol %d for data processing\n", i)
		fmt.Fprintf(&builder, "protocol DataProcessor%d {\n", i)
		builder.WriteString("    // Method signature for processing\n")
		builder.WriteString("    func process(data: [String: Any]) -> [String: Any]\n")
		builder.WriteString("}\n\n")
		fmt.Fprintf(&builder, "// Implementation %d of data processor\n", i)
		fmt.Fprintf(&builder, "struct ConcreteProcessor%d: DataProcessor%d {\n", i, i)
		builder.WriteString("    // Process method implementation\n")
		builder.WriteString("    func process(data: [String: Any]) -> [String: Any] {\n")
		builder.WriteString("        // Input validation\n")
		builder.WriteString("        guard !data.isEmpty else {\n")
		builder.WriteString("            return [:] // Return empty dict for invalid input\n")
		builder.WriteString("        }\n")
		builder.WriteString("        \n")
		builder.WriteString("        // Create processed result\n")
		builder.WriteString("        var result = data // Copy input\n")
		builder.WriteString("        result[\"processed\"] = true // Mark as processed\n")
		fmt.Fprintf(&builder, "        result[\"processorId\"] = %d // Add processor ID\n", i)
		builder.WriteString("        return result // Return result\n")
		builder.WriteString("    }\n")
		builder.WriteString("}\n\n")
	}

	return builder.String()
}

func generateCppCode(targetLines int) string {
	var builder strings.Builder

	builder.WriteString("// C++ high-performance data processing\n")
	builder.WriteString("/* Template-based generic programming\n   for maximum performance and flexibility */\n")
	builder.WriteString("#include <iostream>\n")
	builder.WriteString("#include <vector>\n")
	builder.WriteString("#include <unordered_map>\n")
	builder.WriteString("#include <memory>\n\n")

	functionsToGenerate := targetLines / 16
	for i := 0; i < functionsToGenerate; i++ {
		fmt.Fprintf(&builder, "// Template class %d for generic data processing\n", i)
		fmt.Fprintf(&builder, "template<typename T>\n")
		fmt.Fprintf(&builder, "class DataProcessor%d {\n", i)
		builder.WriteString("private:\n")
		builder.WriteString("    std::vector<T> data_; // Internal data storage\n")
		builder.WriteString("    \n")
		builder.WriteString("public:\n")
		builder.WriteString("    // Constructor with initializer list\n")
		fmt.Fprintf(&builder, "    explicit DataProcessor%d(const std::vector<T>& data) : data_(data) {}\n", i)
		builder.WriteString("    \n")
		builder.WriteString("    // Method to process data with move semantics\n")
		builder.WriteString("    std::vector<T> process() && {\n")
		builder.WriteString("        // Process each element\n")
		builder.WriteString("        for (auto& item : data_) {\n")
		builder.WriteString("            // Apply transformation (placeholder)\n")
		builder.WriteString("            item = std::move(item); // Move optimization\n")
		builder.WriteString("        }\n")
		builder.WriteString("        return std::move(data_); // Return processed data\n")
		builder.WriteString("    }\n")
		builder.WriteString("};\n\n")
	}

	return builder.String()
}

func generateCSharpCode(targetLines int) string {
	var builder strings.Builder

	builder.WriteString("// C# enterprise data processing service\n")
	builder.WriteString("/* .NET Framework compatible implementation\n   with LINQ and async/await patterns */\n")
	builder.WriteString("using System;\n")
	builder.WriteString("using System.Collections.Generic;\n")
	builder.WriteString("using System.Linq;\n")
	builder.WriteString("using System.Threading.Tasks;\n\n")

	functionsToGenerate := targetLines / 14
	for i := 0; i < functionsToGenerate; i++ {
		fmt.Fprintf(&builder, "// Service class %d for async data processing\n", i)
		fmt.Fprintf(&builder, "public class DataService%d\n", i)
		builder.WriteString("{\n")
		builder.WriteString("    // Private readonly field for data storage\n")
		builder.WriteString("    private readonly List<Dictionary<string, object>> _dataStore;\n\n")
		builder.WriteString("    // Constructor with dependency injection\n")
		fmt.Fprintf(&builder, "    public DataService%d()\n", i)
		builder.WriteString("    {\n")
		builder.WriteString("        _dataStore = new List<Dictionary<string, object>>(); // Initialize storage\n")
		builder.WriteString("    }\n\n")
		builder.WriteString("    // Async method for data processing\n")
		builder.WriteString("    public async Task<Dictionary<string, object>> ProcessDataAsync(Dictionary<string, object> input)\n")
		builder.WriteString("    {\n")
		builder.WriteString("        // Input validation with null check\n")
		builder.WriteString("        if (input == null || !input.Any())\n")
		builder.WriteString("        {\n")
		builder.WriteString("            throw new ArgumentException(\"Input cannot be null or empty\"); // Validation error\n")
		builder.WriteString("        }\n\n")
		builder.WriteString("        // Simulate async processing\n")
		builder.WriteString("        await Task.Delay(1); // Non-blocking delay\n\n")
		builder.WriteString("        // Return processed result using LINQ\n")
		fmt.Fprintf(&builder, "        return input.ToDictionary(kvp => kvp.Key, kvp => kvp.Value)\n")
		builder.WriteString("            .Concat(new[] { new KeyValuePair<string, object>(\"processed\", true) })\n")
		fmt.Fprintf(&builder, "            .Concat(new[] { new KeyValuePair<string, object>(\"serviceId\", %d) })\n", i)
		builder.WriteString("            .ToDictionary(kvp => kvp.Key, kvp => kvp.Value); // Final result\n")
		builder.WriteString("    }\n")
		builder.WriteString("}\n\n")
	}

	return builder.String()
}

func generateBashCode(targetLines int) string {
	var builder strings.Builder

	builder.WriteString("#!/bin/bash\n")
	builder.WriteString("# Bash script for system data processing\n")
	builder.WriteString("# Handles file operations and system integration\n\n")
	builder.WriteString("set -euo pipefail # Strict error handling\n\n")

	functionsToGenerate := targetLines / 8
	for i := 0; i < functionsToGenerate; i++ {
		fmt.Fprintf(&builder, "# Function %d: processes system data\n", i)
		fmt.Fprintf(&builder, "process_data_%d() {\n", i)
		builder.WriteString("    local input_file=\"$1\" # Input file parameter\n")
		builder.WriteString("    local output_file=\"$2\" # Output file parameter\n")
		builder.WriteString("    \n")
		builder.WriteString("    # Validate input parameters\n")
		builder.WriteString("    if [[ ! -f \"$input_file\" ]]; then\n")
		builder.WriteString("        echo \"Error: Input file does not exist\" >&2 # Error to stderr\n")
		builder.WriteString("        return 1\n")
		builder.WriteString("    fi\n")
		builder.WriteString("    \n")
		builder.WriteString("    # Process file content\n")
		fmt.Fprintf(&builder, "    sed 's/old/new/g' \"$input_file\" > \"$output_file\" # Transform and save\n")
		fmt.Fprintf(&builder, "    echo \"Processed by function %d\" >> \"$output_file\" # Add signature\n", i)
		builder.WriteString("}\n\n")
	}

	return builder.String()
}

func generateCSSCode(targetLines int) string {
	var builder strings.Builder

	builder.WriteString("/* CSS stylesheet for modern web applications */\n")
	builder.WriteString("/* Responsive design with mobile-first approach */\n\n")

	rulesetsToGenerate := targetLines / 6
	for i := 0; i < rulesetsToGenerate; i++ {
		fmt.Fprintf(&builder, "/* Component %d: styling for data display */\n", i)
		fmt.Fprintf(&builder, ".data-component-%d {\n", i)
		builder.WriteString("    /* Layout properties */\n")
		builder.WriteString("    display: flex;\n")
		builder.WriteString("    flex-direction: column; /* Vertical layout */\n")
		builder.WriteString("    padding: 1rem; /* Internal spacing */\n")
		builder.WriteString("    margin: 0.5rem; /* External spacing */\n")
		fmt.Fprintf(&builder, "    background-color: hsl(%d, 70%%, 95%%); /* Dynamic color */\n", (i*30)%360)
		builder.WriteString("    border-radius: 0.5rem; /* Rounded corners */\n")
		builder.WriteString("    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); /* Subtle shadow */\n")
		builder.WriteString("}\n\n")
		fmt.Fprintf(&builder, "/* Responsive behavior for component %d */\n", i)
		fmt.Fprintf(&builder, "@media (min-width: 768px) {\n")
		fmt.Fprintf(&builder, "    .data-component-%d {\n", i)
		builder.WriteString("        /* Tablet and desktop styles */\n")
		builder.WriteString("        flex-direction: row; /* Horizontal layout */\n")
		builder.WriteString("        padding: 2rem; /* Increased spacing */\n")
		builder.WriteString("    }\n")
		builder.WriteString("}\n\n")
	}

	return builder.String()
}
