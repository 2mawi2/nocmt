package processor

import (
	"regexp"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/python"
)

type PythonSingleProcessor struct {
	*SingleLineCoreProcessor
}


func isPythonSingleLineCommentNode(node *sitter.Node, sourceText string) bool {
	if node.Type() == "comment" {
		commentText := sourceText[node.StartByte():node.EndByte()]
		return strings.HasPrefix(strings.TrimSpace(commentText), "#")
	}
	return false
}


var pythonSingleLineDirectiveRegex = regexp.MustCompile(`(?:\s|^)#\s*(noqa|type:|pylint:|flake8:|mypy:|yapf:|isort:|ruff:|fmt:\s*off|fmt:\s*on)`)


func checkPythonSingleLineDirective(line string) bool {
	trimmedLine := strings.TrimSpace(line)
	if strings.HasPrefix(trimmedLine, "#!") {
		return true
	}
	return pythonSingleLineDirectiveRegex.MatchString(line)
}

func NewPythonSingleProcessor(preserveDirectivesFlag bool) *PythonSingleProcessor {
	singleLineCore := NewSingleLineCoreProcessor(
		"python",
		python.GetLanguage(),
		isPythonSingleLineCommentNode,
		checkPythonSingleLineDirective,
		nil, 
	).WithPreserveDirectives(preserveDirectivesFlag)

	return &PythonSingleProcessor{
		SingleLineCoreProcessor: singleLineCore,
	}
}
