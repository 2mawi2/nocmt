package walker

import (
	"os"
	"path/filepath"
)

type Walker struct{}

type FileProcessor func(path string) error

func (w *Walker) Walk(rootPath string, processor FileProcessor) error {
	checker, err := NewHierarchicalGitIgnoreChecker(rootPath)
	if err != nil {
		return err
	}

	return filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}

			if checker.IsIgnored(path) {
				return filepath.SkipDir
			}

			return nil
		}

		if checker.IsIgnored(path) {
			return nil
		}

		return processor(path)
	})
}