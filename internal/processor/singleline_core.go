package processor

import (
	"context"
	"fmt"
	"strings"

	"nocmt/internal/config"

	sitter "github.com/smacker/go-tree-sitter"
)

type SingleLineCoreProcessor struct {
	langName                string
	lang                    *sitter.Language
	preserveDirectives      bool
	isDirective             func(string) bool
	isSingleLineCommentNode func(node *sitter.Node, sourceText string) bool
	postProcess             func(source string, preserveDirectives bool) (string, error)
	commentConfig           *config.Config
	keepBlankRuns           bool
}

func NewSingleLineCoreProcessor(
	name string,
	lang *sitter.Language,
	isSLCommentNodeFunc func(node *sitter.Node, sourceText string) bool,
	isDirectiveFunc func(string) bool,
	postFunc func(source string, preserveDirectives bool) (string, error),
) *SingleLineCoreProcessor {
	return &SingleLineCoreProcessor{
		langName:                name,
		lang:                    lang,
		isSingleLineCommentNode: isSLCommentNodeFunc,
		isDirective:             isDirectiveFunc,
		postProcess:             postFunc,
		keepBlankRuns:           false,
	}
}

func (p *SingleLineCoreProcessor) WithPreserveDirectives(preserve bool) *SingleLineCoreProcessor {
	p.preserveDirectives = preserve
	return p
}

func (p *SingleLineCoreProcessor) GetLanguageName() string {
	return p.langName
}

func (p *SingleLineCoreProcessor) PreserveDirectives() bool {
	return p.preserveDirectives
}

func (p *SingleLineCoreProcessor) SetCommentConfig(cfg *config.Config) {
	p.commentConfig = cfg
}

func (p *SingleLineCoreProcessor) PreserveBlankRuns() *SingleLineCoreProcessor {
	p.keepBlankRuns = true
	return p
}

func (p *SingleLineCoreProcessor) StripComments(source string) (string, error) {
	if p.lang == nil {
		return source, fmt.Errorf("language %s not configured for Tree-sitter based comment removal", p.langName)
	}

	parser := parsers.Get(p.lang)
	defer parsers.Put(p.lang, parser)

	tree, err := parser.ParseCtx(context.Background(), nil, []byte(source))
	if err != nil {
		return source, fmt.Errorf("failed to parse source for %s: %w", p.langName, err)
	}
	if tree == nil || tree.RootNode() == nil || tree.RootNode().HasError() {
		return source, fmt.Errorf("tree-sitter parsing error for %s, comments not stripped", p.langName)
	}
	defer tree.Close()

	rootNode := tree.RootNode()
	var rangesToModify []CommentRange

	sourceLines := splitIntoLines(source)
	lineStartPositions := calculateLinePositions(sourceLines)

	Walk(rootNode, func(node *sitter.Node) bool {

		if !p.isSingleLineCommentNode(node, source) {
			return true
		}

		commentContent := source[node.StartByte():node.EndByte()]

		// Check user-configured ignore patterns (e.g., "TODO", "FIXME")
		if p.commentConfig != nil && p.commentConfig.ShouldIgnoreComment(commentContent) {
			return true
		}

		if p.preserveDirectives && p.isDirective != nil && p.isDirective(commentContent) {
			return true
		}

		startByte := int(node.StartByte())
		commentLineIndex := -1

		for i, lsp := range lineStartPositions {
			lineContentEndByte := lsp + len(sourceLines[i])
			if startByte >= lsp && startByte <= lineContentEndByte {
				commentLineIndex = i
				break
			}
		}
		if commentLineIndex == -1 && len(sourceLines) > 0 && startByte >= lineStartPositions[len(sourceLines)-1] {
			commentLineIndex = len(sourceLines) - 1
		}

		if commentLineIndex != -1 {
			lineToExamine := sourceLines[commentLineIndex]
			lineStartByteInSource := lineStartPositions[commentLineIndex]

			relativeCommentStartInLine := startByte - lineStartByteInSource

			isEffectivelyFullLine := true
			for k := 0; k < relativeCommentStartInLine; k++ {
				if lineToExamine[k] != ' ' && lineToExamine[k] != '\t' {
					isEffectivelyFullLine = false
					break
				}
			}

			var currentRangeToModify CommentRange
			if isEffectivelyFullLine {

				currentRangeToModify.StartByte = uint32(lineStartByteInSource)
				currentRangeToModify.EndByte = uint32(lineStartByteInSource + len(lineToExamine))

				if commentLineIndex < len(sourceLines)-1 {
					currentRangeToModify.EndByte++
				} else if strings.HasSuffix(source, "\n") && currentRangeToModify.EndByte == uint32(len(source)-1) {

					currentRangeToModify.EndByte++
				}

				currentRangeToModify.Content = ""
			} else {

				adjustedRemoveStart := node.StartByte()
				idx := relativeCommentStartInLine - 1
				for idx >= 0 && (lineToExamine[idx] == ' ' || lineToExamine[idx] == '\t') {
					adjustedRemoveStart = uint32(lineStartByteInSource + idx)
					idx--
				}
				currentRangeToModify.StartByte = adjustedRemoveStart
				currentRangeToModify.EndByte = uint32(lineStartByteInSource + len(lineToExamine))
				currentRangeToModify.Content = ""
			}

			rangesToModify = append(rangesToModify, currentRangeToModify)
		}

		return false
	})

	cleaned := source

	for i := range rangesToModify {
		for j := i + 1; j < len(rangesToModify); j++ {
			if rangesToModify[i].StartByte < rangesToModify[j].StartByte {
				rangesToModify[i], rangesToModify[j] = rangesToModify[j], rangesToModify[i]
			}
		}
	}

	resultBytes := []byte(cleaned)
	for _, r := range rangesToModify {
		start := int(r.StartByte)
		end := int(r.EndByte)

		if start >= len(resultBytes) || end > len(resultBytes) || start > end {
			continue
		}

		if len(r.Content) > 0 {
			resultBytes = append(resultBytes[:start], append([]byte(r.Content), resultBytes[end:]...)...)
		} else {
			resultBytes = append(resultBytes[:start], resultBytes[end:]...)
		}
	}

	cleaned = string(resultBytes)

	if len(rangesToModify) == 0 {
		return source, nil
	}

	if p.postProcess != nil {
		var errPostProcess error
		cleaned, errPostProcess = p.postProcess(cleaned, p.preserveDirectives)
		if errPostProcess != nil {
			return source, errPostProcess
		}
	}

	if p.keepBlankRuns {
		cleaned = normalizeTextKeepBlankRuns(cleaned)
	} else {
		cleaned = normalizeText(cleaned)
	}
	return cleaned, nil
}
