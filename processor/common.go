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

	return findCommentNodes(rootNode, source), nil
}

func parseCode(parser *sitter.Parser, source string) ([]CommentRange, error) {
	return ParseCode(parser, source)
}

func findCommentNodes(node *sitter.Node, source string) []CommentRange {
	var ranges []CommentRange

	if node.Type() == "comment" || node.Type() == "line_comment" || node.Type() == "block_comment" {
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

	sourceLines := strings.Split(source, "\n")
	lineStartOffsets := make([]int, len(sourceLines))

	offset := 0
	for i := range sourceLines {
		lineStartOffsets[i] = offset
		offset += len(sourceLines[i]) + 1
	}

	commentOnlyLines := make(map[int]bool)

	for _, r := range ranges {
		startLine := -1
		endLine := -1

		for i, startOffset := range lineStartOffsets {
			endOffset := startOffset + len(sourceLines[i])
			if int(r.StartByte) >= startOffset && int(r.StartByte) <= endOffset {
				startLine = i
				break
			}
		}

		for i, startOffset := range lineStartOffsets {
			endOffset := startOffset + len(sourceLines[i])
			if int(r.EndByte) >= startOffset && int(r.EndByte) <= endOffset {
				endLine = i
				break
			}
		}

		if startLine != -1 && endLine != -1 {
			for line := startLine; line <= endLine; line++ {
				lineStart := lineStartOffsets[line]
				lineContent := sourceLines[line]

				trimmed := strings.TrimSpace(lineContent)
				if trimmed != "" {
					commentStart := int(r.StartByte) - lineStart
					commentEnd := int(r.EndByte) - lineStart

					if commentStart < 0 {
						commentStart = 0
					}
					if commentEnd > len(lineContent) {
						commentEnd = len(lineContent)
					}

					switch line {
					case startLine:
						beforeComment := strings.TrimSpace(lineContent[:commentStart])
						if beforeComment == "" && line == endLine {
							afterComment := ""
							if commentEnd < len(lineContent) {
								afterComment = strings.TrimSpace(lineContent[commentEnd:])
							}
							if afterComment == "" {
								commentOnlyLines[line] = true
							}
						} else if beforeComment == "" && line != endLine {
							commentOnlyLines[line] = true
						}
					case endLine:
						afterComment := ""
						if commentEnd < len(lineContent) {
							afterComment = strings.TrimSpace(lineContent[commentEnd:])
						}
						if afterComment == "" {
							commentOnlyLines[line] = true
						}
					default:
						commentOnlyLines[line] = true
					}
				} else {
					commentOnlyLines[line] = true
				}
			}
		}
	}

	result := source
	for _, r := range ranges {
		if int(r.StartByte) >= len(result) || int(r.EndByte) > len(result) {
			continue
		}

		beforeComment := result[:r.StartByte]
		afterComment := result[r.EndByte:]

		endsWithNewline := false
		if int(r.EndByte) < len(result) && result[r.EndByte-1] == '\n' {
			endsWithNewline = true
		}

		if endsWithNewline {
			result = beforeComment + "\n" + afterComment
		} else {
			result = beforeComment + afterComment
		}
	}

	lines := strings.Split(result, "\n")
	var cleanedLines []string

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		shouldKeepLine := true

		if commentOnlyLines[i] && trimmedLine == "" {
			shouldKeepLine = false
		} else {
			shouldKeepLine = true
		}

		if shouldKeepLine {
			cleanedLines = append(cleanedLines, line)
		}
	}

	return strings.Join(cleanedLines, "\n")
}

func StripCommentsPreserveDirectives(source string, matcher DirectiveMatcher, parser *sitter.Parser) (string, error) {
	return stripCommentsPreserveDirectives(source, matcher, parser)
}

func stripCommentsPreserveDirectives(source string, matcher DirectiveMatcher, parser *sitter.Parser) (string, error) {
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

func (b *BaseProcessor) stripCommentsWithFiltering(source string, parser *sitter.Parser) (string, error) {
	commentRanges, err := parseCode(parser, source)
	if err != nil {
		return "", err
	}

	commentRanges = b.filterCommentRanges(commentRanges)

	return removeComments(source, commentRanges), nil
}