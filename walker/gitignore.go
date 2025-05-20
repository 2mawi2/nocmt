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

type patternInfo struct {
	pattern    string
	isNegated  bool
	isWildcard bool
	dirLevel   string
}

type GitIgnoreChecker struct {
	ignorer  *gitignore.GitIgnore
	rootPath string
}

func NewGitIgnoreChecker(rootPath string) (*GitIgnoreChecker, error) {
	gitignorePath := filepath.Join(rootPath, ".gitignore")
	var ignorer *gitignore.GitIgnore
	if _, err := os.Stat(gitignorePath); err == nil {
		var err2 error
		ignorer, err2 = gitignore.CompileIgnoreFile(gitignorePath)
		if err2 != nil {
			return nil, err2
		}
	} else {
		ignorer = gitignore.CompileIgnoreLines()
	}
	return &GitIgnoreChecker{ignorer: ignorer, rootPath: rootPath}, nil
}

func (g *GitIgnoreChecker) IsIgnored(path string) bool {
	relPath, err := filepath.Rel(g.rootPath, path)
	if err != nil {
		return false
	}
	relPath = filepath.ToSlash(relPath)
	return g.ignorer.MatchesPath(relPath)
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
		if pat := h.parseLine(line, relDir); pat != nil {
			dirPatterns = append(dirPatterns, *pat)
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
	return &patternInfo{pattern: pattern, isNegated: isNegated, isWildcard: isWildcard, dirLevel: relDir}
}

func (h *HierarchicalGitIgnoreChecker) findApplicableDirectories(relPath string) []string {
	var dirs []string
	for dir := range h.gitignoreFiles {
		if dir == "" || strings.HasPrefix(relPath, dir+"/") || relPath == dir {
			dirs = append(dirs, dir)
		}
	}
	sort.Slice(dirs, func(i, j int) bool {
		return strings.Count(dirs[i], "/") > strings.Count(dirs[j], "/")
	})
	return dirs
}

func (h *HierarchicalGitIgnoreChecker) hasExplicitUnignorePattern(path, relPath string) bool {
	filename := filepath.Base(relPath)
	for dirLevel, pats := range h.patterns {
		if dirLevel == "" || strings.HasPrefix(relPath, dirLevel+"/") || relPath == dirLevel {
			for _, pat := range pats {
				if pat.isNegated {
					if pat.pattern == filename {
						return true
					}
					if pat.isWildcard {
						if m, err := filepath.Match(pat.pattern, filename); err == nil && m {
							return true
						}
					}
				}
			}
		}
	}
	return false
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

	dirs := h.findApplicableDirectories(relPath)
	for _, dir := range dirs {
		ignorer := h.gitignoreFiles[dir]
		subPath := relPath
		if dir != "" {
			prefix := dir + "/"
			if strings.HasPrefix(relPath, prefix) {
				subPath = relPath[len(prefix):]
			} else if relPath == dir {
				subPath = "."
			} else {
				continue
			}
		}
		if subPath == "" {
			continue
		}
		matched, pattern := ignorer.MatchesPathHow(subPath)
		if matched && pattern != nil {
			if pattern.Negate {
				return false
			}
			return true
		}
	}
	return h.defaultIgnorer.MatchesPath(relPath)
}
