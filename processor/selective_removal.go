package processor

import (
	"fmt"
	"nocmt/config"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/bash"
	"github.com/smacker/go-tree-sitter/csharp"
	"github.com/smacker/go-tree-sitter/css"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/rust"
	ts "github.com/smacker/go-tree-sitter/typescript/typescript"
)

func GetParserForProcessor(proc LanguageProcessor) *sitter.Parser {
	parser := sitter.NewParser()

	var language *sitter.Language

	switch proc.GetLanguageName() {
	case "go":
		language = golang.GetLanguage()
	case "javascript":
		language = javascript.GetLanguage()
	case "typescript":
		language = ts.GetLanguage()
	case "python":
		language = python.GetLanguage()
	case "bash":
		language = bash.GetLanguage()
	case "csharp":
		language = csharp.GetLanguage()
	case "rust":
		language = rust.GetLanguage()
	case "css":
		language = css.GetLanguage()
	default:
		return nil
	}

	parser.SetLanguage(language)
	return parser
}

func ParseCodeForCommentRanges(parser *sitter.Parser, source string) ([]CommentRange, error) {
	return ParseCode(parser, source)
}

func FindCommentLineNumbers(content string, comment CommentRange) (startLine, endLine int) {
	lines := strings.Split(content, "\n")
	lineStartOffsets := make([]int, len(lines))

	offset := 0
	for i := range lines {
		lineStartOffsets[i] = offset
		offset += len(lines[i]) + 1
	}

	startLine = -1
	endLine = -1

	for i, startOffset := range lineStartOffsets {
		endOffset := startOffset + len(lines[i])
		if int(comment.StartByte) >= startOffset && int(comment.StartByte) <= endOffset {
			startLine = i
		}
		if int(comment.EndByte-1) >= startOffset && int(comment.EndByte-1) <= endOffset {
			endLine = i
			break
		}
	}

	return startLine + 1, endLine + 1
}

func IsDirective(proc LanguageProcessor, comment string) bool {
	if !proc.PreserveDirectives() {
		return false
	}

	switch proc.GetLanguageName() {
	case "go":
		return checkGoDirective(comment)
	case "javascript":
		return isJSDirective(comment)
	case "typescript":
		return isJSDirective(comment)
	case "python":
		return isPythonDirective(comment)
	case "bash":
		return isBashDirective(comment)
	case "csharp":
		return isCSharpDirective(comment)
	case "rust":
		return isRustDirectiveSelective(comment)
	case "css":
		return isCSSDirective(comment)
	default:
		return false
	}
}

func CommentOverlapsModifiedLines(commentStartLine, endLine int, modifiedLines map[int]bool) bool {
	for line := commentStartLine; line <= endLine; line++ {
		if modifiedLines[line] {
			return true
		}
	}
	return false
}

func FilterCommentsForRemoval(
	commentRanges []CommentRange,
	source string,
	modifiedLines map[int]bool,
	proc LanguageProcessor,
	preserveDirectives bool,
	commentConfig *config.Config,
) []CommentRange {
	var commentsToRemove []CommentRange

	for _, comment := range commentRanges {
		startLine, endLine := FindCommentLineNumbers(source, comment)
		overlapsModified := CommentOverlapsModifiedLines(startLine, endLine, modifiedLines)

		if !overlapsModified {
			continue
		}

		if commentConfig != nil && commentConfig.ShouldIgnoreComment(comment.Content) {
			continue
		}

		if preserveDirectives && IsDirective(proc, comment.Content) {
			continue
		}

		commentsToRemove = append(commentsToRemove, comment)
	}

	return commentsToRemove
}

func SelectivelyStripComments(
	content string,
	filePath string,
	proc LanguageProcessor,
	modifiedLines map[int]bool,
	preserveDirectives bool,
	commentConfig *config.Config,
) (string, error) {
	parser := GetParserForProcessor(proc)
	if parser == nil {
		return "", fmt.Errorf("no tree-sitter parser available for language: %s. Ensure grammar is correctly configured", proc.GetLanguageName())
	}

	commentRanges, err := ParseCodeForCommentRanges(parser, content)
	if err != nil {
		return "", fmt.Errorf("failed to parse code: %w", err)
	}

	commentsToRemove := FilterCommentsForRemoval(commentRanges, content, modifiedLines, proc, preserveDirectives, commentConfig)

	if len(commentsToRemove) == 0 {
		return content, nil
	}

	return RemoveComments(content, commentsToRemove), nil
}

func isPythonDirective(comment string) bool {
	return strings.Contains(comment, "# noqa") ||
		strings.Contains(comment, "# type:") ||
		strings.Contains(comment, "# pragma:") ||
		strings.Contains(comment, "# pylint:")
}

func isBashDirective(comment string) bool {
	return strings.Contains(comment, "# shellcheck") ||
		strings.Contains(comment, "#!")
}

func isRustDirectiveSelective(comment string) bool {
	trimmed := strings.TrimSpace(comment)
	return strings.HasPrefix(trimmed, "#!") || strings.HasPrefix(trimmed, "#[")
}
