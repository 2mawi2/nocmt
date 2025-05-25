package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type CommentConfig struct {
	IgnorePatterns     []string `json:"ignorePatterns"`
	FileIgnorePatterns []string `json:"fileIgnorePatterns"`
}

type Config struct {
	Global               CommentConfig
	Local                CommentConfig
	CLIPatterns          []string
	CLIFilePatterns      []string
	compiledPatterns     []*regexp.Regexp
	compiledFilePatterns []*regexp.Regexp
}

func New() *Config {
	return &Config{}
}

func (c *Config) LoadConfigurations() error {
	homeDir, err := os.UserHomeDir()
	if err == nil {
		globalConfigPath := filepath.Join(homeDir, ".nocmt", "config.json")
		c.Global, _ = loadConfigFile(globalConfigPath)
	}

	localConfigPath := ".nocmt.json"
	c.Local, _ = loadConfigFile(localConfigPath)

	return c.compilePatterns()
}

func (c *Config) SetCLIPatterns(patterns []string) error {
	c.CLIPatterns = patterns
	return c.compilePatterns()
}

func (c *Config) SetCLIFilePatterns(patterns []string) error {
	c.CLIFilePatterns = patterns
	return c.compilePatterns()
}

func (c *Config) ShouldIgnoreComment(comment string) bool {
	for _, pattern := range c.compiledPatterns {
		if pattern.MatchString(comment) {
			return true
		}
	}
	return false
}

func (c *Config) ShouldIgnoreFile(filename string) bool {
	if len(c.compiledFilePatterns) == 0 {
		return false
	}

	filename = filepath.ToSlash(filename)

	for _, pattern := range c.compiledFilePatterns {
		if pattern.MatchString(filename) {
			return true
		}
	}

	if filepath.IsAbs(filename) {
		cwd, err := os.Getwd()
		if err == nil {
			if rel, err := filepath.Rel(cwd, filename); err == nil {
				relPath := filepath.ToSlash(rel)
				for _, pattern := range c.compiledFilePatterns {
					if pattern.MatchString(relPath) {
						return true
					}
				}
			}
		}

		parts := strings.Split(filename, "/")
		for i := 0; i < len(parts); i++ {
			tailPath := filepath.ToSlash(strings.Join(parts[len(parts)-i-1:], "/"))
			for _, pattern := range c.compiledFilePatterns {
				if pattern.MatchString(tailPath) {
					return true
				}
			}
		}
	}

	return false
}

func (c *Config) compilePatterns() error {
	c.compiledPatterns = nil
	c.compiledFilePatterns = nil

	allPatterns := append([]string{}, c.Global.IgnorePatterns...)
	allPatterns = append(allPatterns, c.Local.IgnorePatterns...)
	allPatterns = append(allPatterns, c.CLIPatterns...)

	for _, pattern := range allPatterns {
		compiled, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid pattern '%s': %w", pattern, err)
		}
		c.compiledPatterns = append(c.compiledPatterns, compiled)
	}

	allFilePatterns := append([]string{}, c.Global.FileIgnorePatterns...)
	allFilePatterns = append(allFilePatterns, c.Local.FileIgnorePatterns...)
	allFilePatterns = append(allFilePatterns, c.CLIFilePatterns...)

	for _, pattern := range allFilePatterns {
		compiled, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid file pattern '%s': %w", pattern, err)
		}
		c.compiledFilePatterns = append(c.compiledFilePatterns, compiled)
	}

	return nil
}

func (c *Config) AddIgnorePattern(pattern string) error {
	_, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern '%s': %w", pattern, err)
	}

	c.Local.IgnorePatterns = append(c.Local.IgnorePatterns, pattern)

	return c.SaveLocalConfig()
}

func (c *Config) AddFileIgnorePattern(pattern string) error {
	_, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid file pattern '%s': %w", pattern, err)
	}

	c.Local.FileIgnorePatterns = append(c.Local.FileIgnorePatterns, pattern)

	return c.SaveLocalConfig()
}

func (c *Config) AddGlobalFileIgnorePattern(pattern string) error {
	_, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid file pattern '%s': %w", pattern, err)
	}

	c.Global.FileIgnorePatterns = append(c.Global.FileIgnorePatterns, pattern)

	return c.SaveGlobalConfig()
}

func (c *Config) SaveLocalConfig() error {
	return saveConfigFile(".nocmt.json", c.Local)
}

func (c *Config) SaveGlobalConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot find home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".nocmt")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("cannot create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config.json")
	return saveConfigFile(configPath, c.Global)
}

func (c *Config) AddGlobalIgnorePattern(pattern string) error {
	_, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern '%s': %w", pattern, err)
	}

	c.Global.IgnorePatterns = append(c.Global.IgnorePatterns, pattern)

	return c.SaveGlobalConfig()
}

func loadConfigFile(path string) (CommentConfig, error) {
	config := CommentConfig{}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("error parsing config file %s: %w", path, err)
	}

	return config, nil
}

func saveConfigFile(path string, config CommentConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing config: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}
