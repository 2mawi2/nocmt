package walker

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	gitignore "github.com/sabhiram/go-gitignore"
)

var DefaultIgnorePatterns = []string{
	"package-lock.json",
	"Thumbs.db",
	"desktop.ini",
	".DS_Store",
	"ehthumbs.db",
	"*.swp",
	"*~",
	"._*",
	".Spotlight-V100",
	".Trashes",
	"Icon?",
	".idea/",
	".vscode/",
	"*.sublime-*",
	".project",
	".settings/",
	"*.iml",
	".eclipse",
	"*.suo",
	"*.user",
	"*.sln.docstates",
	"__pycache__/",
	"*.py[cod]",
	"*$py.class",
	"*.so",
	"build/",
	"dist/",
	"*.egg-info/",
	"node_modules/",
	"venv/",
	"env/",
	".chroma/",
	"*.jar",
	"target/",
	"bin/",
	"out/",
	".git/",
	".svn/",
	".hg/",
	"CVS/",
	".gitignore",
	".gitattributes",
	"*.log",
	"*.sqlite",
	"*.db",
	"npm-debug.log*",
	"yarn-debug.log*",
	"yarn-error.log*",
	".env",
	".env.*",
	"*.cfg",
	"*.conf",
	".htaccess",
	"tmp/",
	"temp/",
	"cache/",
	".cache/",
	"*.tmp",
	"*.temp",
	"*.bak",
	"*.backup",
	"*_backup",
	"*.old",
}

var CommonDirectories = []string{
	"node_modules",
	"venv",
	".git",
	"__pycache__",
	".chroma",
	".idea",
	".vscode",
	"dist",
	"build",
	"target",
}

var CommonFiles = []string{
	".ds_store",
	"thumbs.db",
	".gitignore",
	"package-lock.json",
}

var CommonExtensions = []string{
	".pyc",
	".pyo",
	".so",
	".o",
	".class",
	".tmp",
	".swp",
}

type GitIgnoreChecker struct {
	ignorer  *gitignore.GitIgnore
	rootPath string
}

func NewGitIgnoreChecker(rootPath string) (*GitIgnoreChecker, error) {
	gitignorePath := filepath.Join(rootPath, ".gitignore")
	var ignorer *gitignore.GitIgnore

	if _, err := os.Stat(gitignorePath); err == nil {
		var err error
		ignorer, err = gitignore.CompileIgnoreFile(gitignorePath)
		if err != nil {
			return nil, err
		}
	} else {
		ignorer = gitignore.CompileIgnoreLines([]string{}...)
	}

	return &GitIgnoreChecker{
		ignorer:  ignorer,
		rootPath: rootPath,
	}, nil
}

func (g *GitIgnoreChecker) IsIgnored(path string) bool {
	relPath, err := filepath.Rel(g.rootPath, path)
	if err != nil {
		return false
	}

	relPath = filepath.ToSlash(relPath)
	return g.ignorer.MatchesPath(relPath)
}

type gitIgnoreEntry struct {
	dirLevel string
	pattern  *gitignore.IgnorePattern
}

type patternInfo struct {
	pattern    string
	isNegated  bool
	isWildcard bool
	dirLevel   string
}

type HierarchicalGitIgnoreChecker struct {
	gitignoreFiles map[string]*gitignore.GitIgnore
	rootPath       string
	patterns       map[string][]patternInfo
	defaultIgnorer *gitignore.GitIgnore
}

func NewHierarchicalGitIgnoreChecker(rootPath string) (*HierarchicalGitIgnoreChecker, error) {
	checker := &HierarchicalGitIgnoreChecker{
		gitignoreFiles: make(map[string]*gitignore.GitIgnore),
		rootPath:       rootPath,
		patterns:       make(map[string][]patternInfo),
		defaultIgnorer: gitignore.CompileIgnoreLines(DefaultIgnorePatterns...),
	}

	if err := checker.findAllGitignoreFiles(rootPath); err != nil {
		return nil, err
	}

	return checker, nil
}

func (h *HierarchicalGitIgnoreChecker) findAllGitignoreFiles(rootPath string) error {
	return filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && d.Name() == ".gitignore" {
			dir := filepath.Dir(path)
			relDir, err := h.getRelativeDir(dir)
			if err != nil {
				return err
			}

			if err := h.processGitignoreFile(path, relDir); err != nil {
				return err
			}
		}

		return nil
	})
}

func (h *HierarchicalGitIgnoreChecker) getRelativeDir(dir string) (string, error) {
	relDir, err := filepath.Rel(h.rootPath, dir)
	if err != nil {
		return "", err
	}

	relDir = filepath.ToSlash(relDir)
	if relDir == "." {
		return "", nil
	}

	return relDir, nil
}

func (h *HierarchicalGitIgnoreChecker) processGitignoreFile(path, relDir string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var dirPatterns []patternInfo

	for _, line := range lines {
		pattern := h.parseLine(line, relDir)
		if pattern != nil {
			dirPatterns = append(dirPatterns, *pattern)
		}
	}

	h.patterns[relDir] = dirPatterns

	ignorer, err := gitignore.CompileIgnoreFile(path)
	if err != nil {
		return err
	}

	h.gitignoreFiles[relDir] = ignorer
	return nil
}

func (h *HierarchicalGitIgnoreChecker) parseLine(line, relDir string) *patternInfo {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return nil
	}

	isNegated := strings.HasPrefix(line, "!")
	pattern := line
	if isNegated {
		pattern = line[1:]
	}

	isWildcard := strings.Contains(pattern, "*")

	return &patternInfo{
		pattern:    pattern,
		isNegated:  isNegated,
		isWildcard: isWildcard,
		dirLevel:   relDir,
	}
}

func (h *HierarchicalGitIgnoreChecker) ShouldIgnoreByDefault(path string) bool {
	if h.isCommonFileOrExtension(path) {
		return true
	}

	if h.isInCommonDirectory(path) {
		return true
	}

	return h.matchesDefaultPattern(path)
}

func (h *HierarchicalGitIgnoreChecker) isCommonFileOrExtension(path string) bool {
	filename := filepath.Base(path)
	filenameCase := strings.ToLower(filename)

	for _, commonFile := range CommonFiles {
		if strings.ToLower(commonFile) == filenameCase {
			return true
		}
	}

	for _, ext := range CommonExtensions {
		if strings.HasSuffix(filenameCase, ext) {
			return true
		}
	}

	return false
}

func (h *HierarchicalGitIgnoreChecker) isInCommonDirectory(path string) bool {
	relPath, err := filepath.Rel(h.rootPath, path)
	if err != nil {
		return false
	}

	relPath = filepath.ToSlash(relPath)
	pathParts := strings.Split(relPath, "/")

	for _, part := range pathParts {
		partCase := strings.ToLower(part)
		for _, commonDir := range CommonDirectories {
			if partCase == strings.ToLower(commonDir) {
				return true
			}
		}
	}

	return false
}

func (h *HierarchicalGitIgnoreChecker) matchesDefaultPattern(path string) bool {
	relPath, err := filepath.Rel(h.rootPath, path)
	if err != nil {
		return false
	}

	relPath = filepath.ToSlash(relPath)
	return h.defaultIgnorer.MatchesPath(relPath)
}

func (h *HierarchicalGitIgnoreChecker) IsIgnored(path string) bool {
	relPath, err := filepath.Rel(h.rootPath, path)
	if err != nil {
		return false
	}

	relPath = filepath.ToSlash(relPath)

	if h.hasExplicitUnignorePattern(path, relPath) {
		return false
	}

	applicableDirs := h.findApplicableDirectories(relPath)

	if h.hasSpecialLogUnignorePattern(relPath) {
		return false
	}

	unignoreMatches := h.findUnignoreMatches(relPath, applicableDirs)

	if h.hasIgnorePattern(relPath, applicableDirs, unignoreMatches) {
		return true
	}

	if len(unignoreMatches) == 0 && h.ShouldIgnoreByDefault(path) {
		return true
	}

	return false
}

func (h *HierarchicalGitIgnoreChecker) hasExplicitUnignorePattern(path, relPath string) bool {
	filename := filepath.Base(relPath)

	for dirLevel, patterns := range h.patterns {
		if dirLevel == "" || strings.HasPrefix(relPath, dirLevel+"/") || relPath == dirLevel {
			for _, pat := range patterns {
				if pat.isNegated && (pat.pattern == filename || strings.TrimPrefix(pat.pattern, "!") == filename) {
					return true
				}
			}
		}
	}

	return false
}

func (h *HierarchicalGitIgnoreChecker) findApplicableDirectories(relPath string) []string {
	var applicableDirs []string
	for dirLevel := range h.gitignoreFiles {
		if dirLevel == "" || strings.HasPrefix(relPath, dirLevel+"/") || relPath == dirLevel {
			applicableDirs = append(applicableDirs, dirLevel)
		}
	}

	sort.Slice(applicableDirs, func(i, j int) bool {
		return strings.Count(applicableDirs[i], "/") > strings.Count(applicableDirs[j], "/")
	})

	return applicableDirs
}

func (h *HierarchicalGitIgnoreChecker) hasSpecialLogUnignorePattern(relPath string) bool {
	if strings.HasSuffix(relPath, ".log") {
		dirPath := filepath.Dir(relPath)
		if dirPath == "." {
			dirPath = ""
		}

		if dirPath != "" {
			if patterns, ok := h.patterns[dirPath]; ok {
				for _, pat := range patterns {
					if pat.isNegated && (pat.pattern == "*.log" || pat.pattern == strings.TrimPrefix(pat.pattern, "!")) {
						return true
					}
				}
			}
		}
	}

	return false
}

func (h *HierarchicalGitIgnoreChecker) findUnignoreMatches(relPath string, applicableDirs []string) []gitIgnoreEntry {
	var unignoreMatches []gitIgnoreEntry

	for _, dirLevel := range applicableDirs {
		ignorer := h.gitignoreFiles[dirLevel]
		dirRelPath := h.getRelativePathForDir(relPath, dirLevel)
		if dirRelPath == "" {
			continue
		}

		matches, pattern := ignorer.MatchesPathHow(dirRelPath)

		if matches && pattern != nil && pattern.Negate {
			unignoreMatches = append(unignoreMatches, gitIgnoreEntry{
				dirLevel: dirLevel,
				pattern:  pattern,
			})

			if dirLevel == applicableDirs[0] {
				return unignoreMatches
			}
		}
	}

	return unignoreMatches
}

func (h *HierarchicalGitIgnoreChecker) getRelativePathForDir(relPath, dirLevel string) string {
	if dirLevel == "" {
		return relPath
	}

	prefix := dirLevel + "/"
	if strings.HasPrefix(relPath, prefix) {
		return relPath[len(prefix):]
	}

	if relPath == dirLevel {
		return "."
	}

	return ""
}

func (h *HierarchicalGitIgnoreChecker) hasIgnorePattern(relPath string, applicableDirs []string, unignoreMatches []gitIgnoreEntry) bool {
	for _, dirLevel := range applicableDirs {
		ignorer := h.gitignoreFiles[dirLevel]
		dirRelPath := h.getRelativePathForDir(relPath, dirLevel)
		if dirRelPath == "" {
			continue
		}

		matches, pattern := ignorer.MatchesPathHow(dirRelPath)

		if matches && pattern != nil && !pattern.Negate {
			return !h.isOverriddenByMoreSpecificUnignore(dirLevel, unignoreMatches)
		}
	}

	return false
}

func (h *HierarchicalGitIgnoreChecker) isOverriddenByMoreSpecificUnignore(dirLevel string, unignoreMatches []gitIgnoreEntry) bool {
	for _, unignore := range unignoreMatches {
		if strings.Count(unignore.dirLevel, "/") > strings.Count(dirLevel, "/") {
			return true
		}
	}
	return false
}
