package walker

import (
	"fmt"
	"os"
	"path/filepath"

	"nocmt/config"
	"nocmt/processor"
)

type ProcessorConfig struct {
	PreserveDirectives bool
	DryRun             bool
	Verbose            bool
	Force              bool
	CommentConfig      *config.Config
}

type ProcessorIntegration struct {
	factory        *processor.ProcessorFactory
	config         ProcessorConfig
	processedCount int
	skippedCount   int
	errorCount     int
}

func NewProcessorIntegration(config ProcessorConfig) *ProcessorIntegration {
	factory := processor.NewProcessorFactory()
	factory.SetPreserveDirectives(config.PreserveDirectives)
	if config.CommentConfig != nil {
		factory.SetCommentConfig(config.CommentConfig)
	}

	return &ProcessorIntegration{
		factory: factory,
		config:  config,
	}
}

func (p *ProcessorIntegration) ProcessRepository(rootPath string) error {
	walker := &Walker{}
	return walker.Walk(rootPath, p.processFile)
}

func (p *ProcessorIntegration) GetStats() (processed, skipped, errors int) {
	return p.processedCount, p.skippedCount, p.errorCount
}

func (p *ProcessorIntegration) processFile(path string) error {
	if p.config.CommentConfig != nil && p.config.CommentConfig.ShouldIgnoreFile(path) {
		p.skippedCount++
		if p.config.Verbose {
			fmt.Printf("Skipping %s: matches file ignore pattern\n", path)
		}
		return nil
	}

	proc, err := p.factory.GetProcessorByExtension(filepath.Base(path))
	if err != nil {
		p.skippedCount++
		if p.config.Verbose {
			fmt.Printf("Skipping %s: %v\n", path, err)
		}
		return nil
	}

	if proc.GetLanguageName() == "go" && p.config.PreserveDirectives {
		proc = processor.NewGoProcessor(true)
		if p.config.CommentConfig != nil {
			proc.SetCommentConfig(p.config.CommentConfig)
		}
	}

	content, err := os.ReadFile(path)
	if err != nil {
		p.errorCount++
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	strippedContent, err := proc.StripComments(string(content))
	if err != nil {
		p.errorCount++
		return fmt.Errorf("failed to process %s: %w", path, err)
	}

	if strippedContent == string(content) {
		p.skippedCount++
		if p.config.Verbose {
			fmt.Printf("No changes needed for %s\n", path)
		}
		return nil
	}

	if p.config.Verbose {
		fmt.Printf("Processing %s\n", path)
	}

	if !p.config.DryRun {
		err = os.WriteFile(path, []byte(strippedContent), 0644)
		if err != nil {
			p.errorCount++
			return fmt.Errorf("failed to write %s: %w", path, err)
		}
	}

	p.processedCount++
	return nil
}