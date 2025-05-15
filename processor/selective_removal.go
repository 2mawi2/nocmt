package processor

import (
	"fmt"
	"nocmt/config"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/bash"
	"github.com/smacker/go-tree-sitter/css"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/rust"
)

func GetParserForProcessor(proc LanguageProcessor) *sitter.Parser {
	parser := sitter.NewParser()

	var language *sitter.Language

	switch proc.GetLanguageName() {
	case "go":
		language = golang.GetLanguage()
	case "javascript", "typescript":
		language = javascript.GetLanguage()
	case "java":
		language = java.GetLanguage()
	case "python":
		language = python.GetLanguage()
	case "csharp":
		// TODO: Add proper C# language support
		return nil
	case "rust":
		language = rust.GetLanguage()
	case "kotlin":
		// TODO: Add proper Kotlin language support
		return nil
	case "bash":
		language = bash.GetLanguage()
	case "swift":
		// TODO: Add proper Swift language support
		return nil
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
		return isGoDirective(comment)
	case "javascript", "typescript":
		return isJSDirective(comment)
	case "python":
		return isPythonDirective(comment)
	case "bash":
		return isBashDirective(comment)
	case "java":
		return isJavaDirective(comment)
	case "csharp":
		return isCSharpDirective(comment)
	case "rust":
		return isRustDirective(comment)
	case "kotlin":
		return isKotlinDirective(comment)
	case "swift":
		return isSwiftDirective(comment)
	case "css":
		return isCSSDirective(comment)
	default:
		return false
	}
}

func CommentOverlapsModifiedLines(commentStartLine, commentEndLine int, modifiedLines map[int]bool) bool {
	for line := commentStartLine; line <= commentEndLine; line++ {
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
		return "", fmt.Errorf("no tree-sitter parser available for language: %s", proc.GetLanguageName())
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

func isGoDirective(comment string) bool {
	return strings.Contains(comment, "//go:") || strings.Contains(comment, "// go:")
}

func isJSDirective(comment string) bool {
	return strings.Contains(comment, "@ts-") ||
		strings.Contains(comment, "// @") ||
		strings.Contains(comment, "/* @") ||
		strings.Contains(comment, "/*global") ||
		strings.Contains(comment, "/*eslint")
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

func isJavaDirective(comment string) bool {
	return strings.Contains(comment, "@") ||
		strings.Contains(comment, "// FIXME:") ||
		strings.Contains(comment, "// TODO:") ||
		strings.Contains(comment, "// XXX:")
}

func isCSharpDirective(comment string) bool {
	return strings.Contains(comment, "#pragma") ||
		strings.Contains(comment, "#region") ||
		strings.Contains(comment, "#endregion") ||
		strings.Contains(comment, "#nullable")
}

func isRustDirective(comment string) bool {
	return strings.Contains(comment, "#!") ||
		strings.Contains(comment, "//!")
}

func isKotlinDirective(comment string) bool {
	return strings.Contains(comment, "@") ||
		strings.Contains(comment, "@file:")
}

func isSwiftDirective(comment string) bool {
	return strings.Contains(comment, "// MARK:") ||
		strings.Contains(comment, "// TODO:") ||
		strings.Contains(comment, "// FIXME:")
}

func isCSSDirective(comment string) bool {
	return strings.Contains(comment, "/*!")
}
