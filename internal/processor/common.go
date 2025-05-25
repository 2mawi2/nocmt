package processor

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"nocmt/internal/config"

	sitter "github.com/smacker/go-tree-sitter"
)

type DirectiveMatcher func(line string) bool

type CommentRange struct {
	StartByte, EndByte uint32
	Content            string
}

type BaseProcessor struct {
	commentConfig *config.Config
}

func (b *BaseProcessor) SetCommentConfig(cfg *config.Config) {
	b.commentConfig = cfg
}

func (b *BaseProcessor) ShouldIgnoreComment(comment string) bool {
	if b.commentConfig == nil {
		return false
	}
	return b.commentConfig.ShouldIgnoreComment(comment)
}

func ParseCode(parser *sitter.Parser, source string) ([]CommentRange, error) {
	sourceBytes := []byte(source)
	tree, err := parser.ParseCtx(context.Background(), nil, sourceBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse source code: %w", err)
	}
	if tree == nil {
		return nil, fmt.Errorf("failed to parse source code")
	}
	defer tree.Close()

	rootNode := tree.RootNode()
	if rootNode == nil {
		return nil, fmt.Errorf("failed to get root node")
	}

	if rootNode.HasError() {
		return nil, fmt.Errorf("syntax error in source code (rootNode.HasError() is true)")
	}

	return findCommentNodes(rootNode, source), nil
}

func Walk(node *sitter.Node, callback func(*sitter.Node) bool) {
	if !callback(node) {
		return
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		Walk(node.Child(i), callback)
	}
}

func parseCode(parser *sitter.Parser, source string) ([]CommentRange, error) {
	return ParseCode(parser, source)
}

func findCommentNodes(node *sitter.Node, source string) []CommentRange {
	var ranges []CommentRange

	switch node.Type() {
	case "comment",
		"line_comment",
		"block_comment",
		"documentation_comment",
		"doc_comment":
		ranges = append(ranges, CommentRange{
			StartByte: node.StartByte(),
			EndByte:   node.EndByte(),
			Content:   source[node.StartByte():node.EndByte()],
		})
	}

	childCount := int(node.ChildCount())
	for i := 0; i < childCount; i++ {
		child := node.Child(i)
		if !child.IsNull() {
			ranges = append(ranges, findCommentNodes(child, source)...)
		}
	}

	return ranges
}

func (b *BaseProcessor) filterCommentRanges(ranges []CommentRange) []CommentRange {
	if b.commentConfig == nil {
		return ranges
	}

	var filteredRanges []CommentRange
	for _, r := range ranges {
		if !b.ShouldIgnoreComment(r.Content) {
			filteredRanges = append(filteredRanges, r)
		}
	}

	return filteredRanges
}

func RemoveComments(source string, ranges []CommentRange) string {
	return removeComments(source, ranges)
}

func removeComments(source string, ranges []CommentRange) string {
	if len(ranges) == 0 {
		return source
	}

	for i := range ranges {
		for j := i + 1; j < len(ranges); j++ {
			if ranges[i].StartByte < ranges[j].StartByte {
				ranges[i], ranges[j] = ranges[j], ranges[i]
			}
		}
	}

	resultBytes := []byte(source)
	for _, r := range ranges {
		start := int(r.StartByte)
		end := int(r.EndByte)

		if start > len(resultBytes) || end > len(resultBytes) || start > end {
			continue
		}

		resultBytes = append(resultBytes[:start], resultBytes[end:]...)
	}

	return string(resultBytes)
}

func StripCommentsPreserveDirectives(source string, matcher DirectiveMatcher, parser *sitter.Parser) (string, error) {
	lines := strings.Split(source, "\n")
	directiveLines := make(map[int]string)

	for i, line := range lines {
		if matcher(line) {
			directiveLines[i] = line
		}
	}

	commentRanges, err := parseCode(parser, source)
	if err != nil {
		return "", err
	}

	stripped := removeComments(source, commentRanges)

	strippedLines := strings.Split(stripped, "\n")
	for i, directive := range directiveLines {
		if i < len(strippedLines) {
			strippedLines[i] = directive
		}
	}

	return strings.Join(strippedLines, "\n"), nil
}

func PreserveOriginalTrailingNewline(original, cleaned string) string {
	originalHadTrailing := strings.HasSuffix(original, "\n")
	if !originalHadTrailing && strings.HasSuffix(cleaned, "\n") {
		return strings.TrimSuffix(cleaned, "\n")
	}
	if originalHadTrailing && !strings.HasSuffix(cleaned, "\n") {
		return cleaned + "\n"
	}
	return cleaned
}

func splitIntoLines(sourceCode string) []string {
	if shouldUseFastPathForSmallStrings(sourceCode) {
		return strings.Split(sourceCode, "\n")
	}

	return splitLargeStringIntoLinesManually(sourceCode)
}

func shouldUseFastPathForSmallStrings(sourceCode string) bool {
	const smallStringThreshold = 100
	return len(sourceCode) < smallStringThreshold
}

func splitLargeStringIntoLinesManually(sourceCode string) []string {
	const estimatedCharsPerLine = 50
	estimatedLineCount := len(sourceCode) / estimatedCharsPerLine
	lines := make([]string, 0, estimatedLineCount)

	currentLineStart := 0

	for currentPos := 0; currentPos < len(sourceCode); currentPos++ {
		if isUnixNewline(sourceCode, currentPos) {
			lines = appendLineFromRange(lines, sourceCode, currentLineStart, currentPos)
			currentLineStart = currentPos + 1
		} else if isWindowsNewline(sourceCode, currentPos) {
			lines = appendLineFromRange(lines, sourceCode, currentLineStart, currentPos)
			currentPos++ 
			currentLineStart = currentPos + 1
		}
	}

	return appendRemainingContentAsLastLine(lines, sourceCode, currentLineStart)
}

func isUnixNewline(sourceCode string, position int) bool {
	return sourceCode[position] == '\n'
}

func isWindowsNewline(sourceCode string, position int) bool {
	return sourceCode[position] == '\r' &&
		position+1 < len(sourceCode) &&
		sourceCode[position+1] == '\n'
}

func appendLineFromRange(lines []string, sourceCode string, startPos, endPos int) []string {
	return append(lines, sourceCode[startPos:endPos])
}

func appendRemainingContentAsLastLine(lines []string, sourceCode string, lastLineStart int) []string {
	if hasRemainingContent(sourceCode, lastLineStart) {
		lines = append(lines, sourceCode[lastLineStart:])
	}
	return lines
}

func hasRemainingContent(sourceCode string, startPosition int) bool {
	return startPosition < len(sourceCode)
}

func calculateLinePositions(lines []string) []int {
	positions := make([]int, len(lines))
	pos := 0
	for i, line := range lines {
		positions[i] = pos
		pos += len(line) + 1
	}
	return positions
}

func normalizeText(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	trimmedLines := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmedLines = append(trimmedLines, strings.TrimRight(line, " \t"))
	}
	s = strings.Join(trimmedLines, "\n")

	s = regexp.MustCompile(`\n{3,}`).ReplaceAllString(s, "\n\n")

	s = regexp.MustCompile(`^\n+`).ReplaceAllString(s, "")

	s = regexp.MustCompile(`\n+$`).ReplaceAllString(s, "")

	if s != "" {
		s += "\n"
	}
	return s
}

func normalizeTextKeepBlankRuns(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")

	lines := strings.Split(s, "\n")
	for i, ln := range lines {
		lines[i] = strings.TrimRight(ln, " \t")
	}
	s = strings.Join(lines, "\n")

	if s != "" && !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	return s
}

type ParserPoolType struct {
	sync.Mutex
	parsers map[*sitter.Language]*sync.Pool
}

var parsers = &ParserPoolType{
	parsers: make(map[*sitter.Language]*sync.Pool),
}

func (p *ParserPoolType) Get(lang *sitter.Language) *sitter.Parser {
	p.Lock()
	defer p.Unlock()

	pool, exists := p.parsers[lang]
	if !exists {
		pool = &sync.Pool{
			New: func() interface{} {
				parser := sitter.NewParser()
				parser.SetLanguage(lang)
				return parser
			},
		}
		p.parsers[lang] = pool
	}

	return pool.Get().(*sitter.Parser)
}

func (p *ParserPoolType) Put(lang *sitter.Language, parser *sitter.Parser) {
	p.Lock()
	defer p.Unlock()

	if pool, exists := p.parsers[lang]; exists {
		pool.Put(parser)
	}
}
