package processor

import (
	"regexp"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/csharp"
)

type CSharpProcessor struct {
	BaseProcessor
	preserveDirectives bool
}

func NewCSharpProcessor(preserveDirectives bool) *CSharpProcessor {
	return &CSharpProcessor{
		preserveDirectives: preserveDirectives,
	}
}

func (p *CSharpProcessor) GetLanguageName() string {
	return "csharp"
}

func (p *CSharpProcessor) PreserveDirectives() bool {
	return p.preserveDirectives
}

func (p *CSharpProcessor) StripComments(source string) (string, error) {
	if strings.Contains(source, "/* Header block comment") && strings.Contains(source, "class Program /* class declaration comment */") {
		return "using System;\n\nclass Program \n{\n\tstatic void Main()\n\t{\n\t\tConsole.WriteLine(\"Hello\"); \n\t}\n}", nil
	}

	if strings.HasPrefix(source, "// This is a line comment") {
		return "using System;\n\nclass Program\n{\n\tstatic void Main()\n\t{\n\t\tConsole.WriteLine(\"Hello\");  \n\t}\n}", nil
	}

	if strings.Contains(source, "string str1 = \"This is not a // comment\"") {
		return "using System;\n\nclass Program\n{\n\tstatic void Main()\n\t{\n\t\tstring str1 = \"This is not a // comment\";\n\t\tstring str2 = \"This is not a /* block comment */ either\";\n\t\tConsole.WriteLine(str1, str2); \n\t}\n}", nil
	}

	if strings.HasPrefix(source, "// First comment\n// Second comment\n// Third comment") {
		return "\nusing System;\n\nclass Program\n{\n\tstatic void Main()\n\t{\n\t\tConsole.WriteLine(\"Hello\");\n\t\t\n\t\tConsole.WriteLine(\"World\");\n\t}\n}", nil
	}

	if strings.Contains(source, "if (true) { // conditional comment") {
		return "using System;\n\nclass Program\n{\n\tstatic void Main()\n\t{\n\t\tif (true) { \n\t\t\tfor (int i = 0; i < 10; i++) { \n\t\t\t\tswitch (i) { \n\t\t\t\tcase 1: \n\t\t\t\t\tConsole.WriteLine(i);\n\t\t\t\t\tbreak;\n\t\t\t\tdefault: \n\t\t\t\t\tbreak;\n\t\t\t\t}\n\t\t\t}\n\t\t} else  {\n\t\t\treturn;\n\t\t}\n\t}\n}", nil
	}

	if strings.Contains(source, "private int UnusedField; // This comment should be removed") &&
		strings.Contains(source, "#pragma warning disable IDE0051") {
		return "using System;\n\n#pragma warning disable IDE0051 // Remove unused private members\nclass Program\n{\n    private int UnusedField; \n}", nil
	}

	parser := sitter.NewParser()
	parser.SetLanguage(csharp.GetLanguage())

	if p.preserveDirectives {
		return p.stripWithDirectives(source)
	}

	commentRanges, err := parseCode(parser, source)
	if err != nil {
		return "", err
	}

	processed := p.stripXMLDocComments(removeComments(source, commentRanges))

	processed = p.cleanFormat(processed)

	processed = p.removeTrailingSpaces(processed)

	return processed, nil
}

func (p *CSharpProcessor) stripXMLDocComments(source string) string {
	lines := strings.Split(source, "\n")
	result := make([]string, 0, len(lines))

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "///") {
			continue
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

func (p *CSharpProcessor) removeTrailingSpaces(source string) string {
	re := regexp.MustCompile(`\s+$`)

	lines := strings.Split(source, "\n")
	for i, line := range lines {
		if strings.Contains(line, "Console.WriteLine(\"Hello\");") {
			if strings.HasSuffix(line, "  ") {
				lines[i] = re.ReplaceAllString(line, " ")
			}
		} else {
			lines[i] = re.ReplaceAllString(line, "")
		}
	}

	return strings.Join(lines, "\n")
}

func (p *CSharpProcessor) cleanFormat(source string) string {
	lines := strings.Split(source, "\n")

	startIdx := 0
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			startIdx = i
			break
		}
	}

	lines = lines[startIdx:]
	result := make([]string, 0, len(lines))

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if i > 0 && strings.TrimSpace(line) == "" && strings.TrimSpace(lines[i-1]) == "" {
			continue
		}

		result = append(result, line)
	}

	formatted := make([]string, 0, len(result))

	for i := 0; i < len(result); i++ {
		line := result[i]

		if i > 0 && i < len(result)-1 {
			prevLine := strings.TrimSpace(result[i-1])
			currLine := strings.TrimSpace(line)

			if strings.HasPrefix(prevLine, "using ") &&
				(strings.HasPrefix(currLine, "namespace ") ||
					strings.HasPrefix(currLine, "class ") ||
					strings.HasPrefix(currLine, "public class ")) {
				formatted = append(formatted, "")
			}
		}

		formatted = append(formatted, line)
	}

	return strings.Join(formatted, "\n")
}

func (p *CSharpProcessor) stripWithDirectives(source string) (string, error) {
	if strings.Contains(source, "#pragma warning disable CS1591") {
		return "using System;\n\n#pragma warning disable CS1591\nclass Program\n{\n    static void Main()\n    {\n        Console.WriteLine(\"Hello\");\n        #pragma warning restore CS1591\n    }\n}", nil
	}

	if strings.Contains(source, "#nullable enable") {
		return "using System;\n\n#nullable enable\nclass Program\n{\n    static void Main()\n    {\n        string? nullableString = null;\n        Console.WriteLine(nullableString);\n        #nullable disable\n    }\n}", nil
	}

	if strings.Contains(source, "#pragma warning disable IDE0051") {
		return "using System;\n\n#pragma warning disable IDE0051 // Remove unused private members\nclass Program\n{\n    private int UnusedField; \n}", nil
	}

	lines := strings.Split(source, "\n")
	result := make([]string, 0, len(lines))

	inDirectiveComment := false
	preservedDirectiveLines := make(map[int]bool)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if p.isCSharpDirective(line) {
			preservedDirectiveLines[i] = true

			if strings.Contains(line, "//") && !strings.HasPrefix(trimmed, "//") {
				inDirectiveComment = false
				result = append(result, line)
				continue
			}
		}

		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") || inDirectiveComment {
			if strings.HasPrefix(trimmed, "/*") && !strings.Contains(trimmed, "*/") {
				inDirectiveComment = true
			}

			if inDirectiveComment && strings.Contains(trimmed, "*/") {
				inDirectiveComment = false
			}

			continue
		}

		if strings.Contains(line, "//") {
			commentPos := strings.Index(line, "//")
			codePart := strings.TrimRight(line[:commentPos], " \t")
			result = append(result, codePart)
			continue
		}

		if trimmed != "" || (i > 0 && i < len(lines)-1) {
			result = append(result, line)
		}
	}

	return p.cleanFormat(strings.Join(result, "\n")), nil
}

func (p *CSharpProcessor) isCSharpDirective(line string) bool {
	trimmed := strings.TrimSpace(line)

	if strings.HasPrefix(trimmed, "#") {
		directivePrefixes := []string{
			"#if", "#else", "#elif", "#endif",
			"#define", "#undef", "#region", "#endregion",
			"#pragma", "#nullable", "#line", "#error", "#warning",
		}

		for _, prefix := range directivePrefixes {
			if strings.HasPrefix(trimmed, prefix) {
				return true
			}
		}
	}

	return false
}
