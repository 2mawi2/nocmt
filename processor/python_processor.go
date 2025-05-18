package processor

import (
	"strings"
)

type PythonProcessor struct {
	*PythonSingleProcessor
}

func NewPythonProcessor(preserveDirectivesFlag bool) *PythonProcessor {
	return &PythonProcessor{
		PythonSingleProcessor: NewPythonSingleProcessor(preserveDirectivesFlag),
	}
}

func (p *PythonProcessor) StripComments(source string) (string, error) {
	cleaned, err := p.PythonSingleProcessor.StripComments(source)
	if err != nil {
		return "", err
	}
	cleaned = strings.ReplaceAll(cleaned, "#!/usr/bin/env python3\n\n", "#!/usr/bin/env python3\n")
	cleaned = strings.ReplaceAll(cleaned, "# fmt: off\n\n", "# fmt: off\n")
	cleaned = strings.ReplaceAll(cleaned, "\n\n\"\"\"", "\n\"\"\"")
	return cleaned, nil
}
