package processor

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/swift"
)

type SwiftProcessor struct {
	*SingleLineCoreProcessor
}

func isSwiftSingleLineCommentNode(node *sitter.Node, sourceText string) bool {
	if node.Type() == "comment" {
		commentText := sourceText[node.StartByte():node.EndByte()]
		trimmed := strings.TrimSpace(commentText)
		if strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "///") {
			return true
		}
	}
	return false
}

func isSwiftDirective(line string) bool {
	trimmed := strings.TrimSpace(line)

	
	if strings.HasPrefix(trimmed, "// @") {
		
		afterPrefix := trimmed[4:] 
		validAttributes := []string{
			"available", "objc", "objcMembers", "nonobjc", "NSManaged", "IBAction",
			"IBOutlet", "IBDesignable", "IBInspectable", "GKInspectable",
			"discardableResult", "warn_unqualified_access", "autoclosure",
			"escaping", "inlinable", "usableFromInline", "frozen", "propertyWrapper",
		}
		for _, attr := range validAttributes {
			if strings.HasPrefix(afterPrefix, attr) {
				nextChar := len(attr)
				if nextChar >= len(afterPrefix) || (afterPrefix[nextChar] < 'a' || afterPrefix[nextChar] > 'z') && (afterPrefix[nextChar] < 'A' || afterPrefix[nextChar] > 'Z') && (afterPrefix[nextChar] < '0' || afterPrefix[nextChar] > '9') {
					return true
				}
			}
		}
	}

	return strings.HasPrefix(trimmed, "// TODO:") ||
		strings.HasPrefix(trimmed, "// FIXME:") ||
		strings.HasPrefix(trimmed, "// MARK:") ||
		strings.HasPrefix(trimmed, "// WARNING:") ||
		strings.HasPrefix(trimmed, "// NOTE:") ||
		strings.Contains(trimmed, "swiftlint:") ||
		strings.Contains(trimmed, "sourcery:")
}

func NewSwiftProcessor(preserveDirectives bool) *SwiftProcessor {
	singleLineCore := NewSingleLineCoreProcessor(
		"swift",
		swift.GetLanguage(),
		isSwiftSingleLineCommentNode,
		isSwiftDirective,
		nil,
	).WithPreserveDirectives(preserveDirectives)

	return &SwiftProcessor{
		SingleLineCoreProcessor: singleLineCore,
	}
}

func (p *SwiftProcessor) GetLanguageName() string {
	return "swift"
}

func (p *SwiftProcessor) PreserveDirectives() bool {
	return p.preserveDirectives
}

func (p *SwiftProcessor) StripComments(source string) (string, error) {
	cleaned, err := p.SingleLineCoreProcessor.StripComments(source)
	if err != nil {
		return "", err
	}
	return PreserveOriginalTrailingNewline(source, cleaned), nil
}
