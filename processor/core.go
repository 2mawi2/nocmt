package processor

import (
	"regexp"
	"strings"
	"sync"

	"nocmt/config"

	sitter "github.com/smacker/go-tree-sitter"
)

type CoreProcessor struct {
	langName           string
	lang               *sitter.Language
	preserveDirectives bool
	isDirective        func(string) bool
	postProcess        func(source string, commentRanges []CommentRange, preserveDirectives bool) (string, error)
	commentConfig      *config.Config
	keepBlankRuns      bool
}

func NewCoreProcessor(
	name string,
	lang *sitter.Language,
	isDir func(string) bool,
	post func(source string, commentRanges []CommentRange, preserveDirectives bool) (string, error),
) *CoreProcessor {
	return &CoreProcessor{
		langName:      name,
		lang:          lang,
		isDirective:   isDir,
		postProcess:   post,
		keepBlankRuns: false,
	}
}

func (p *CoreProcessor) WithPreserveDirectives(preserve bool) *CoreProcessor {
	p.preserveDirectives = preserve
	return p
}

func (p *CoreProcessor) GetLanguageName() string {
	return p.langName
}

func (p *CoreProcessor) PreserveDirectives() bool {
	return p.preserveDirectives
}

func (p *CoreProcessor) SetCommentConfig(cfg *config.Config) {
	p.commentConfig = cfg
}

func (p *CoreProcessor) PreserveBlankRuns() *CoreProcessor {
	p.keepBlankRuns = true
	return p
}

func (p *CoreProcessor) StripComments(source string) (string, error) {
	parser := parsers.Get(p.lang)
	defer parsers.Put(p.lang, parser)

	commentRanges, err := parseCode(parser, source)
	if err != nil {
		return "", err
	}

	if p.commentConfig != nil {
		commentRanges = filterConfigIgnores(source, commentRanges, p.commentConfig)
	}

	activeCommentRanges := commentRanges
	if p.isDirective != nil && p.preserveDirectives {
		var nonDirectiveCommentRanges []CommentRange
		lines := splitIntoLines(source)
		lineStartPositions := calculateLinePositions(lines)
		directiveMap := make(map[int]bool)
		for i, line := range lines {
			if p.isDirective(line) {
				directiveMap[i] = true
			}
		}

		for _, r := range commentRanges {
			isAssociatedWithDirective := false
			startIdx, endIdx := getCommentLineIndices(source, r, lineStartPositions, lines)
			for i := startIdx; i <= endIdx; i++ {
				if directiveMap[i] {
					isAssociatedWithDirective = true
					break
				}
			}
			if !isAssociatedWithDirective {
				nonDirectiveCommentRanges = append(nonDirectiveCommentRanges, r)
			}
		}
		activeCommentRanges = nonDirectiveCommentRanges
	}

	cleaned := removeComments(source, activeCommentRanges)

	var errPostProcess error
	if p.postProcess != nil {
		cleaned, errPostProcess = p.postProcess(cleaned, commentRanges, p.preserveDirectives)
		if errPostProcess != nil {
			return "", errPostProcess
		}
	}
	if p.keepBlankRuns {
		cleaned = normalizeTextKeepBlankRuns(cleaned)
	} else {
		cleaned = normalizeText(cleaned)
	}
	return cleaned, nil
}

func getCommentLineIndices(sourceContent string, r CommentRange, lineStartPositions []int, lines []string) (int, int) {
	startLine := -1
	endLine := -1

	for i, startPos := range lineStartPositions {
		if int(r.StartByte) < startPos+len(lines[i])+1 {
			startLine = i
			break
		}
	}
	if startLine == -1 && len(lines) > 0 {
		if int(r.StartByte) >= lineStartPositions[len(lines)-1]+len(lines[len(lines)-1]) {
			startLine = len(lines) - 1
		}
	}

	for i := startLine; i < len(lines) && i >= 0; i++ {
		if int(r.EndByte) <= lineStartPositions[i]+len(lines[i])+1 {
			endLine = i
			break
		}
	}
	if endLine == -1 && startLine != -1 {
		if int(r.EndByte) > lineStartPositions[len(lines)-1]+len(lines[len(lines)-1]) {
			endLine = len(lines) - 1
		} else {
			for i := len(lines) - 1; i >= startLine; i-- {
				if int(r.EndByte) > lineStartPositions[i] {
					endLine = i
					break
				}
			}
		}
	}

	if startLine == -1 {
		startLine = 0
	}
	if endLine == -1 {
		endLine = len(lines) - 1
	}
	if endLine < startLine {
		endLine = startLine
	}

	return startLine, endLine
}

func filterConfigIgnores(source string, ranges []CommentRange, cfg *config.Config) []CommentRange {
	if len(ranges) == 0 || cfg == nil {
		return ranges
	}

	var filtered []CommentRange
	for _, r := range ranges {
		if !cfg.ShouldIgnoreComment(r.Content) {
			filtered = append(filtered, r)
		}
	}

	return filtered
}

func splitIntoLines(s string) []string {
	return lineRegexp.Split(s, -1)
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

var lineRegexp = regexp.MustCompile(`\r?\n`)

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
