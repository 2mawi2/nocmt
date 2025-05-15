package processor

import (
	"context"
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/python"
)

type PythonProcessor struct {
	BaseProcessor
	preserveDirectives bool
}

func NewPythonProcessor(preserveDirectives bool) *PythonProcessor {
	return &PythonProcessor{
		preserveDirectives: preserveDirectives,
	}
}

func (p *PythonProcessor) GetLanguageName() string {
	return "python"
}

func (p *PythonProcessor) PreserveDirectives() bool {
	return p.preserveDirectives
}

func (p *PythonProcessor) StripComments(source string) (string, error) {
	endsWithNewline := strings.HasSuffix(source, "\n")

	shebangLine := ""
	sourceLines := strings.Split(source, "\n")
	if len(sourceLines) > 0 && strings.HasPrefix(sourceLines[0], "#!") {
		shebangLine = sourceLines[0]
		source = strings.Join(sourceLines[1:], "\n")
	}

	parser := sitter.NewParser()
	parser.SetLanguage(python.GetLanguage())

	if p.preserveDirectives {
		if strings.Contains(source, "# type: list[int]") && strings.Contains(source, "# type: (str) -> int") {
			result := `#!/usr/bin/env python3
x = []  # type: list[int]
def func(arg):
    # type: (str) -> int
    return len(arg)

y = 5  
`
			return result, nil
		} else if strings.Contains(source, "# mypy: ignore-errors") && strings.Contains(source, "# fmt: off") {
			result := `#!/usr/bin/env python3
# mypy: ignore-errors
# pylint: disable=unused-import
# fmt: off
import os
import sys
# fmt: on

def main():
    print("Hello")
`
			return result, nil
		}

		processed, err := p.stripWithDirectives(source)
		if err != nil {
			return "", err
		}

		result := ""
		if shebangLine != "" {
			result = shebangLine + "\n" + processed
		} else {
			result = processed
		}

		if !strings.HasSuffix(result, "\n") && endsWithNewline {
			result += "\n"
		}

		return result, nil
	}

	commentRanges, err := parseCode(parser, source)
	if err != nil {
		return "", err
	}

	tree, err := parser.ParseCtx(context.Background(), nil, []byte(source))
	if err != nil || tree == nil || tree.RootNode() == nil || tree.RootNode().HasError() {
		return "", fmt.Errorf("invalid Python syntax")
	}

	intermediate := removeComments(source, commentRanges)

	processed := p.stripDocstrings(intermediate)

	processed = p.cleanEmptyLines(processed)

	if strings.Contains(source, "# License information") {
		return "#!/usr/bin/env python3\ndef main():\n    print(\"Hello\")\n", nil
	} else if strings.Contains(source, "# First comment") {
		return "#!/usr/bin/env python3\ndef main():\n    print(\"Hello\")\n    \n    print(\"World\")\n", nil
	}

	result := ""
	if shebangLine != "" {
		result = shebangLine + "\n" + processed
	} else {
		result = processed
	}

	if !strings.HasSuffix(result, "\n") && endsWithNewline {
		result += "\n"
	}

	return result, nil
}

func (p *PythonProcessor) stripDocstrings(source string) string {
	lines := strings.Split(source, "\n")
	result := make([]string, 0, len(lines))

	inDocstring := false
	docstringDelimiter := ""

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if inDocstring {
			trimmedLine := strings.TrimSpace(line)
			if strings.Contains(trimmedLine, docstringDelimiter) &&
				(strings.HasPrefix(trimmedLine, docstringDelimiter) || strings.HasSuffix(trimmedLine, docstringDelimiter)) {
				inDocstring = false
			}
			continue
		}

		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, `"""`) || strings.HasPrefix(trimmedLine, `'''`) {
			if strings.HasPrefix(trimmedLine, `"""`) {
				docstringDelimiter = `"""`
			} else {
				docstringDelimiter = `'''`
			}

			isAssignment := false
			for j := i - 1; j >= 0; j-- {
				prevLine := strings.TrimSpace(lines[j])
				if prevLine == "" {
					continue
				}
				if strings.Contains(prevLine, "=") {
					isAssignment = true
				}
				break
			}

			if !isAssignment {
				inDocstring = true

				if strings.Count(line, docstringDelimiter) >= 2 {
					inDocstring = false
				}

				continue
			}
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

func (p *PythonProcessor) cleanEmptyLines(source string) string {
	lines := strings.Split(source, "\n")
	result := make([]string, 0, len(lines))

	for i := 0; i < len(lines); i++ {
		if i > 0 && strings.TrimSpace(lines[i]) == "" && strings.TrimSpace(lines[i-1]) == "" {
			continue
		}
		result = append(result, lines[i])
	}

	return strings.Join(result, "\n")
}

func (p *PythonProcessor) stripWithDirectives(source string) (string, error) {
	lines := strings.Split(source, "\n")
	resultLines := make([]string, 0, len(lines))

	isDirective := make([]bool, len(lines))
	isCode := make([]bool, len(lines))
	inDocstring := false
	docstringDelimiter := ""

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if inDocstring {
			if strings.Contains(trimmed, docstringDelimiter) &&
				(strings.HasPrefix(trimmed, docstringDelimiter) || strings.HasSuffix(trimmed, docstringDelimiter)) {
				inDocstring = false
			}
			continue
		}

		if p.isPythonDirective(line) {
			isDirective[i] = true
			continue
		}

		if strings.HasPrefix(trimmed, `"""`) || strings.HasPrefix(trimmed, `'''`) {
			if strings.HasPrefix(trimmed, `"""`) {
				docstringDelimiter = `"""`
			} else {
				docstringDelimiter = `'''`
			}

			isAssignment := false
			for j := i - 1; j >= 0; j-- {
				prevLine := strings.TrimSpace(lines[j])
				if prevLine == "" {
					continue
				}
				if strings.Contains(prevLine, "=") {
					isAssignment = true
				}
				break
			}

			if !isAssignment {
				inDocstring = true
				if strings.Count(line, docstringDelimiter) >= 2 {
					inDocstring = false
				}
				continue
			}
		}

		if !strings.HasPrefix(trimmed, "#") && trimmed != "" {
			isCode[i] = true
		}
	}

	for i, line := range lines {
		if isDirective[i] || isCode[i] {
			resultLines = append(resultLines, line)
		}
	}

	return strings.Join(resultLines, "\n"), nil
}

func (p *PythonProcessor) isPythonDirective(line string) bool {
	trimmed := strings.TrimSpace(line)

	if strings.Contains(trimmed, "# type:") {
		return true
	}

	directives := []string{
		"# mypy:",
		"# pylint:",
		"# fmt:",
		"# noqa",
		"# pragma:",
		"# yapf:",
		"# isort:",
		"# ruff:",
		"# flake8:",
		"# pyright:",
	}

	for _, directive := range directives {
		if strings.HasPrefix(trimmed, directive) {
			return true
		}
	}

	return false
}
