package processor

import (
	"fmt"
	"strings"

	"nocmt/config"
	// TODO: Add swift grammar once a reliable one is found or created
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
	factory.Register(NewJavaProcessor(false))
	factory.Register(NewPythonProcessor(false))
	factory.Register(NewCSharpProcessor(false))
	factory.Register(NewRustProcessor(false))
	factory.Register(NewKotlinProcessor(false))
	factory.Register(NewBashProcessor(false))
	factory.Register(NewSwiftProcessor(false))
	factory.Register(NewCSSProcessor(false))

	factory.RegisterConstructor("go", func(preserveDirectives bool) LanguageProcessor {
		return NewGoProcessor(preserveDirectives)
	})
	factory.RegisterConstructor("javascript", func(preserveDirectives bool) LanguageProcessor {
		return NewJavaScriptProcessor(preserveDirectives)
	})
	factory.RegisterConstructor("typescript", func(preserveDirectives bool) LanguageProcessor {
		return NewTypeScriptProcessor(preserveDirectives)
	})
	factory.RegisterConstructor("java", func(preserveDirectives bool) LanguageProcessor {
		return NewJavaProcessor(preserveDirectives)
	})
	factory.RegisterConstructor("python", func(preserveDirectives bool) LanguageProcessor {
		return NewPythonProcessor(preserveDirectives)
	})
	factory.RegisterConstructor("csharp", func(preserveDirectives bool) LanguageProcessor {
		return NewCSharpProcessor(preserveDirectives)
	})
	factory.RegisterConstructor("rust", func(preserveDirectives bool) LanguageProcessor {
		return NewRustProcessor(preserveDirectives)
	})
	factory.RegisterConstructor("kotlin", func(preserveDirectives bool) LanguageProcessor {
		return NewKotlinProcessor(preserveDirectives)
	})
	factory.RegisterConstructor("bash", func(preserveDirectives bool) LanguageProcessor {
		return NewBashProcessor(preserveDirectives)
	})
	factory.RegisterConstructor("swift", func(preserveDirectives bool) LanguageProcessor {
		return NewSwiftProcessor(preserveDirectives)
	})
	factory.RegisterConstructor("css", func(preserveDirectives bool) LanguageProcessor {
		return NewCSSProcessor(preserveDirectives)
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
	if ok {
		processor := constructor(f.preserveDirectives)
		if f.commentConfig != nil {
			processor.SetCommentConfig(f.commentConfig)
		}
		return processor, nil
	}

	processor, ok := f.processors[language]
	if !ok {
		return nil, fmt.Errorf("no processor available for language: %s", language)
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
		".tsx":   "typescript",
		".java":  "java",
		".py":    "python",
		".pyi":   "python",
		".pyx":   "python",
		".cs":    "csharp",
		".rs":    "rust",
		".kt":    "kotlin",
		".kts":   "kotlin",
		".swift": "swift",
		".sh":    "bash",
		".bash":  "bash",
		".css":   "css",
		".scss":  "css",
		".less":  "css",
	}

	for ext, lang := range extMap {
		if strings.HasSuffix(filename, ext) {
			return f.GetProcessor(lang)
		}
	}

	return nil, fmt.Errorf("no processor available for file: %s", filename)
}

func StripComments(source string) (string, error) {
	processor := NewGoProcessor(false)
	return processor.StripComments(source)
}
