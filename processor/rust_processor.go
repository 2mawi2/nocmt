package processor

import (
	"bytes"
	"context"
	"sort"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/rust"
)

type interval struct{ start, end uint32 }

type RustProcessor struct{ *CoreProcessor }

func NewRustProcessor(preserveDirectivesFlag bool) *RustProcessor {
	core := NewCoreProcessor(
		"rust",
		rust.GetLanguage(),
		isRustDirective,
		postProcessRust, 
	).WithPreserveDirectives(preserveDirectivesFlag)
	return &RustProcessor{CoreProcessor: core}
}



func isRustDirective(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "#[") || strings.HasPrefix(trimmed, "#!")
}




func postProcessRust(src string, _ []CommentRange, preserve bool) (string, error) {
	out, err := stripCommentsRust(src, preserve)
	if err != nil {
		return "", err
	}
	return out, nil
}



func stripCommentsRust(source string, preserveDirectives bool) (string, error) {
	if strings.TrimSpace(source) == "" {
		return "", nil
	}
	

	parser := sitter.NewParser()
	parser.SetLanguage(rust.GetLanguage())
	tree, err := parser.ParseCtx(context.Background(), nil, []byte(source))
	if err != nil {
		return "", err
	}

	

	var gaps []interval
	cursor := sitter.NewTreeCursor(tree.RootNode())

	for {
		node := cursor.CurrentNode()
		if node == nil {
			break
		}

		switch node.Type() {
		case "comment":
			gaps = append(gaps, interval{node.StartByte(), node.EndByte()})
		case "attribute_item", "inner_attribute_item", "outer_attribute_item":
			if !preserveDirectives {
				gaps = append(gaps, interval{node.StartByte(), node.EndByte()})
			}
		}

		if cursor.GoToFirstChild() {
			continue
		}
		for !cursor.GoToNextSibling() {
			if !cursor.GoToParent() {
				goto done
			}
		}
	}
done:

	

	sort.Slice(gaps, func(i, j int) bool { return gaps[i].start < gaps[j].start })
	merged := make([]interval, 0, len(gaps))
	for _, iv := range gaps {
		if n := len(merged); n > 0 && iv.start <= merged[n-1].end {
			if iv.end > merged[n-1].end {
				merged[n-1].end = iv.end
			}
		} else {
			merged = append(merged, iv)
		}
	}

	

	var buf bytes.Buffer
	prev := uint32(0)
	srcBytes := []byte(source)
	for _, iv := range merged {
		buf.Write(srcBytes[prev:iv.start])
		prev = iv.end
	}
	buf.Write(srcBytes[prev:])

	

	lines := strings.Split(buf.String(), "\n")
	out := lines[:0]

	for _, ln := range lines {
		trimmed := strings.TrimRight(ln, " \t\r")
		if !preserveDirectives && isRustDirective(strings.TrimSpace(trimmed)) {
			continue
		}
		if trimmed == "" {
			if len(out) > 0 && out[len(out)-1] == "" {
				continue
			}
			out = append(out, "")
		} else {
			
			leading := 0
			for leading < len(trimmed) && (trimmed[leading] == ' ' || trimmed[leading] == '\t') {
				leading++
			}
			indent := trimmed[:leading]
			rest := trimmed[leading:]
			collapsed := indent + strings.Join(strings.Fields(rest), " ")
			out = append(out, collapsed)
		}
	}
	result := strings.Join(out, "\n")
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	return result, nil
}
