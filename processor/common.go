package processor

import (
	"context"
	"fmt"
	"strings"

	"nocmt/config"

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
