package processor

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
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

func stripPythonDocStrings(source string) (string, error) {
	parser := parsers.Get(python.GetLanguage())
	defer parsers.Put(python.GetLanguage(), parser)

	tree, err := parser.ParseCtx(context.Background(), nil, []byte(source))
	if err != nil {
		return source, fmt.Errorf("failed to parse Python for docstring removal: %w", err)
	}
	if tree.RootNode().HasError() {
		return source, nil
	}

	query, err := sitter.NewQuery([]byte(`
        (module (expression_statement (string) @docstring_module))
        (function_definition body: (block (expression_statement (string) @docstring_func)))
        (class_definition body: (block (expression_statement (string) @docstring_class)))
    `), python.GetLanguage())
	if err != nil {
		return source, fmt.Errorf("failed to create Python docstring query: %w", err)
	}

	cursor := sitter.NewQueryCursor()
	cursor.Exec(query, tree.RootNode())

	var docstringRanges []CommentRange
	processedNodes := make(map[uintptr]bool)

	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			node := capture.Node
			if processedNodes[node.ID()] {
				continue
			}

			expressionStmtNode := node.Parent()
			if expressionStmtNode == nil || expressionStmtNode.Type() != "expression_statement" {
				continue
			}

			bodyNode := expressionStmtNode.Parent()
			if bodyNode == nil {
				continue
			}

			firstRelevantChild := true
			var currentChild *sitter.Node
			childCount := bodyNode.NamedChildCount()
			foundOurNode := false
			for i := 0; i < int(childCount); i++ {
				currentChild = bodyNode.NamedChild(i)
				if currentChild == nil {
					continue
				}
				if currentChild.Type() == "decorator" || currentChild.Type() == "comment" || currentChild.Type() == "type_comment" {
					continue
				}
				if currentChild.ID() == expressionStmtNode.ID() {
					foundOurNode = true
					break
				} else {
					firstRelevantChild = false
					break
				}
			}

			if foundOurNode && firstRelevantChild {
				docstringRanges = append(docstringRanges, CommentRange{
					StartByte: uint32(node.StartByte()),
					EndByte:   uint32(node.EndByte()),
					Content:   source[node.StartByte():node.EndByte()],
				})
				processedNodes[node.ID()] = true
			}
		}
	}
	return removeComments(source, docstringRanges), nil
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
