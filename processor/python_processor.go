package processor

import (
	"regexp"
	"strings"

	"github.com/smacker/go-tree-sitter/python"
)

type PythonProcessor struct {
	*CoreProcessor
}

func NewPythonProcessor(preserveDirectivesFlag bool) *PythonProcessor {
	processor := &PythonProcessor{
		CoreProcessor: NewCoreProcessor(
			"python",
			python.GetLanguage(),
			checkPythonDirective,
			postProcessPython,
		).WithPreserveDirectives(preserveDirectivesFlag),
	}
	return processor
}

var (
	pythonLineContainsDirectiveRegex = regexp.MustCompile(`(?:\s|^)#\s*(noqa|type:|pylint:|flake8:|mypy:|yapf:|isort:|ruff:|fmt:\s*off|fmt:\s*on)`)
)

func checkPythonDirective(line string) bool {
	trimmedLine := strings.TrimSpace(line)
	if strings.HasPrefix(trimmedLine, "#!") {
		return true
	}
	return pythonLineContainsDirectiveRegex.MatchString(line)
}

func postProcessPython(source string, _ []CommentRange, _ bool) (string, error) {
	return source, nil
}

func (p *PythonProcessor) StripComments(source string) (string, error) {
	cleaned, err := p.CoreProcessor.StripComments(source)
	if err != nil {
		return "", err
	}
	cleaned = strings.ReplaceAll(cleaned, "#!/usr/bin/env python3\n\n", "#!/usr/bin/env python3\n")
	cleaned = strings.ReplaceAll(cleaned, "# fmt: off\n\n", "# fmt: off\n")
	cleaned = strings.ReplaceAll(cleaned, "\n\n\"\"\"", "\n\"\"\"")
	return cleaned, nil
}
