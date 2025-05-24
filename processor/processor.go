package processor

import (
	"fmt"
	"strings"

	"nocmt/config"
)

type LanguageProcessor interface {
	StripComments(source string) (string, error)

	GetLanguageName() string

	PreserveDirectives() bool

	SetCommentConfig(cfg *config.Config)
}

type ProcessorFactory struct {
	processors            map[string]LanguageProcessor
	preserveDirectives    bool
	commentConfig         *config.Config
	processorConstructors map[string]func(bool) LanguageProcessor
}

func NewProcessorFactory() *ProcessorFactory {
	factory := &ProcessorFactory{
		processors:            make(map[string]LanguageProcessor),
		processorConstructors: make(map[string]func(bool) LanguageProcessor),
		preserveDirectives:    false,
	}

	factory.Register(NewGoProcessor(false))
	factory.Register(NewJavaScriptProcessor(false))
	factory.Register(NewTypeScriptProcessor(false))
	factory.Register(NewPythonSingleProcessor(false))
	factory.Register(NewCSharpSingleProcessor(false))
	factory.Register(NewRustProcessor(false))
	factory.Register(NewBashProcessor(false))
	factory.Register(NewCSSProcessor(false))
	factory.Register(NewKotlinProcessor(false))
	factory.Register(NewJavaProcessor(false))
	factory.Register(NewSwiftProcessor(false))

	factory.RegisterConstructor("go", func(preserveDirectives bool) LanguageProcessor {
		return NewGoProcessor(preserveDirectives)
	})
	factory.RegisterConstructor("javascript", func(preserveDirectives bool) LanguageProcessor {
		return NewJavaScriptProcessor(preserveDirectives)
	})
	factory.RegisterConstructor("typescript", func(preserveDirectives bool) LanguageProcessor {
		return NewTypeScriptProcessor(preserveDirectives)
	})
	factory.RegisterConstructor("python", func(preserveDirectives bool) LanguageProcessor {
		return NewPythonSingleProcessor(preserveDirectives)
	})
	factory.RegisterConstructor("csharp", func(preserveDirectives bool) LanguageProcessor {
		return NewCSharpSingleProcessor(preserveDirectives)
	})
	factory.RegisterConstructor("rust", func(preserveDirectives bool) LanguageProcessor {
		return NewRustProcessor(preserveDirectives)
	})
	factory.RegisterConstructor("bash", func(preserveDirectives bool) LanguageProcessor {
		return NewBashProcessor(preserveDirectives)
	})
	factory.RegisterConstructor("css", func(preserveDirectives bool) LanguageProcessor {
		return NewCSSProcessor(preserveDirectives)
	})
	factory.RegisterConstructor("kotlin", func(preserveDirectives bool) LanguageProcessor {
		return NewKotlinProcessor(preserveDirectives)
	})
	factory.RegisterConstructor("java", func(preserveDirectives bool) LanguageProcessor {
		return NewJavaProcessor(preserveDirectives)
	})
	factory.RegisterConstructor("swift", func(preserveDirectives bool) LanguageProcessor {
		return NewSwiftProcessor(preserveDirectives)
	})

	return factory
}

func (f *ProcessorFactory) SetPreserveDirectives(preserveDirectives bool) {
	f.preserveDirectives = preserveDirectives
}

func (f *ProcessorFactory) SetCommentConfig(cfg *config.Config) {
	f.commentConfig = cfg
}

func (f *ProcessorFactory) Register(processor LanguageProcessor) {
	f.processors[processor.GetLanguageName()] = processor
}

func (f *ProcessorFactory) RegisterConstructor(language string, constructor func(bool) LanguageProcessor) {
	f.processorConstructors[language] = constructor
}

func (f *ProcessorFactory) GetProcessor(language string) (LanguageProcessor, error) {
	constructor, ok := f.processorConstructors[language]
	var processor LanguageProcessor
	if ok {
		processor = constructor(f.preserveDirectives)
	} else {
		var ok2 bool
		processor, ok2 = f.processors[language]
		if !ok2 {
			return nil, fmt.Errorf("no processor available for language: %s", language)
		}
	}
	if f.commentConfig != nil {
		processor.SetCommentConfig(f.commentConfig)
	}
	return processor, nil
}

func (f *ProcessorFactory) GetProcessorByExtension(filename string) (LanguageProcessor, error) {
	extMap := map[string]string{
		".go":    "go",
		".js":    "javascript",
		".jsx":   "javascript",
		".ts":    "typescript",
		".py":    "python",
		".pyi":   "python",
		".pyx":   "python",
		".cs":    "csharp",
		".rs":    "rust",
		".sh":    "bash",
		".bash":  "bash",
		".css":   "css",
		".scss":  "css",
		".less":  "css",
		".kt":    "kotlin",
		".kts":   "kotlin",
		".java":  "java",
		".swift": "swift",
	}

	for ext, lang := range extMap {
		if strings.HasSuffix(filename, ext) {
			if ext == ".tsx" {
				return &noOpProcessor{}, nil
			}
			return f.GetProcessor(lang)
		}
	}

	return nil, fmt.Errorf("no processor available for file: %s", filename)
}

func StripComments(source string) (string, error) {
	panic("StripComments is deprecated. Use ProcessorFactory to obtain a language-specific processor.")
}

type noOpProcessor struct{}

func (n *noOpProcessor) StripComments(source string) (string, error) { return source, nil }
func (n *noOpProcessor) GetLanguageName() string                     { return "tsx" }
func (n *noOpProcessor) PreserveDirectives() bool                    { return false }
func (n *noOpProcessor) SetCommentConfig(cfg *config.Config)         {}
