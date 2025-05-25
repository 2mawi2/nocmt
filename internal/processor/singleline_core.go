package processor

import (
	"context"
	"fmt"
	"sort"
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

func findLineContainingBytePosition(targetBytePosition int, lineStartPositions []int, sourceCode string) int {
	if len(lineStartPositions) == 0 {
		return -1
	}

	leftBoundary, rightBoundary := 0, len(lineStartPositions)-1

	for leftBoundary <= rightBoundary {
		middleIndex := (leftBoundary + rightBoundary) / 2
		lineStartByte := lineStartPositions[middleIndex]
		lineEndByte := calculateLineEndByte(middleIndex, lineStartPositions, sourceCode)

		if isBytePositionWithinLine(targetBytePosition, lineStartByte, lineEndByte) {
			return middleIndex
		} else if targetBytePosition < lineStartByte {
			rightBoundary = middleIndex - 1
		} else {
			leftBoundary = middleIndex + 1
		}
	}

	return fallbackToLastLineIfExists(lineStartPositions)
}

func calculateLineEndByte(lineIndex int, lineStartPositions []int, sourceCode string) int {
	isLastLine := lineIndex == len(lineStartPositions)-1
	if isLastLine {
		return len(sourceCode)
	}
	return lineStartPositions[lineIndex+1] - 1
}

func isBytePositionWithinLine(targetPosition, lineStart, lineEnd int) bool {
	return targetPosition >= lineStart && targetPosition <= lineEnd
}

func fallbackToLastLineIfExists(lineStartPositions []int) int {
	if len(lineStartPositions) > 0 {
		return len(lineStartPositions) - 1
	}
	return -1
}

func sortCommentRangesInReverseOrderToAvoidOffsetIssues(commentRanges []CommentRange) {
	sort.Slice(commentRanges, func(i, j int) bool {
		return commentRanges[i].StartByte > commentRanges[j].StartByte
	})
}

func createCommentRangeForLine(
	commentNode *sitter.Node,
	lineIndex int,
	sourceLines []string,
	lineStartPositions []int,
	fullSourceCode string,
) CommentRange {
	lineContent := sourceLines[lineIndex]
	lineStartByte := lineStartPositions[lineIndex]
	commentStartByte := int(commentNode.StartByte())
	commentPositionInLine := commentStartByte - lineStartByte

	if isCommentOnOtherwiseEmptyLine(lineContent, commentPositionInLine) {
		return createRangeForFullLineComment(lineIndex, lineContent, lineStartByte, sourceLines, fullSourceCode)
	}

	return createRangeForPartialLineComment(commentNode, lineContent, lineStartByte, commentPositionInLine)
}

func isCommentOnOtherwiseEmptyLine(lineContent string, commentStartPosition int) bool {
	for i := range commentStartPosition {
		if lineContent[i] != ' ' && lineContent[i] != '\t' {
			return false
		}
	}
	return true
}

func createRangeForFullLineComment(
	lineIndex int,
	lineContent string,
	lineStartByte int,
	allSourceLines []string,
	fullSourceCode string,
) CommentRange {
	commentRange := CommentRange{
		StartByte: uint32(lineStartByte),
		EndByte:   uint32(lineStartByte + len(lineContent)),
		Content:   "",
	}

	if shouldIncludeTrailingNewline(lineIndex, allSourceLines, fullSourceCode, commentRange.EndByte) {
		commentRange.EndByte++
	}

	return commentRange
}

func shouldIncludeTrailingNewline(
	lineIndex int,
	allSourceLines []string,
	fullSourceCode string,
	currentEndByte uint32,
) bool {
	isNotLastLine := lineIndex < len(allSourceLines)-1
	if isNotLastLine {
		return true
	}

	isLastLineWithTrailingNewline := strings.HasSuffix(fullSourceCode, "\n") &&
		currentEndByte == uint32(len(fullSourceCode)-1)
	return isLastLineWithTrailingNewline
}

func createRangeForPartialLineComment(
	commentNode *sitter.Node,
	lineContent string,
	lineStartByte int,
	commentPositionInLine int,
) CommentRange {
	adjustedStartByte := findStartOfWhitespaceBeforeComment(
		commentNode.StartByte(),
		lineContent,
		lineStartByte,
		commentPositionInLine,
	)

	return CommentRange{
		StartByte: adjustedStartByte,
		EndByte:   uint32(lineStartByte + len(lineContent)),
		Content:   "",
	}
}

func findStartOfWhitespaceBeforeComment(
	originalCommentStart uint32,
	lineContent string,
	lineStartByte int,
	commentPositionInLine int,
) uint32 {
	adjustedStart := originalCommentStart

	for i := commentPositionInLine - 1; i >= 0; i-- {
		if lineContent[i] == ' ' || lineContent[i] == '\t' {
			adjustedStart = uint32(lineStartByte + i)
		} else {
			break
		}
	}

	return adjustedStart
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

		commentStartByte := int(node.StartByte())
		commentLineIndex := findLineContainingBytePosition(commentStartByte, lineStartPositions, source)

		if commentLineIndex != -1 {
			commentRange := createCommentRangeForLine(
				node,
				commentLineIndex,
				sourceLines,
				lineStartPositions,
				source,
			)
			rangesToModify = append(rangesToModify, commentRange)
		}

		return false
	})

	cleaned := source

	sortCommentRangesInReverseOrderToAvoidOffsetIssues(rangesToModify)

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
