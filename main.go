package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"nocmt/config"
	"nocmt/processor"
	"nocmt/walker"
)

func main() {
	var preserveDirectives bool
	var removeDirectives bool
	var dryRun bool
	var verbose bool
	var force bool
	var ignorePatterns string
	var ignoreFilePatterns string
	var configAdd string
	var configAddGlobal string
	var configAddFileIgnore string
	var configAddFileIgnoreGlobal string
	var staged bool
	var all bool
	var showVersion bool

	flag.BoolVar(&removeDirectives, "remove-directives", false, "Remove compiler directives (preserved by default)")
	flag.BoolVar(&removeDirectives, "r", false, "Remove compiler directives (shorthand)")
	flag.BoolVar(&dryRun, "dry-run", false, "Preview changes without modifying files")
	flag.BoolVar(&dryRun, "d", false, "Preview changes without modifying files (shorthand)")
	flag.BoolVar(&verbose, "verbose", false, "Show detailed output during processing")
	flag.BoolVar(&verbose, "v", false, "Show detailed output (shorthand)")
	flag.BoolVar(&force, "force", false, "Run in non-git directories")
	flag.BoolVar(&force, "f", false, "Run in non-git directories (shorthand)")
	flag.StringVar(&ignorePatterns, "ignore", "", "Comma-separated list of regex patterns to preserve comments")
	flag.StringVar(&ignoreFilePatterns, "ignore-file", "", "Comma-separated list of regex patterns to ignore files")
	flag.StringVar(&configAdd, "add-ignore", "", "Add a regex pattern to the project's ignore list")
	flag.StringVar(&configAddGlobal, "add-ignore-global", "", "Add a regex pattern to the global ignore list")
	flag.StringVar(&configAddFileIgnore, "add-ignore-file", "", "Add a regex pattern to the local file ignore list")
	flag.StringVar(&configAddFileIgnoreGlobal, "add-ignore-file-global", "", "Add a regex pattern to the global file ignore list")
	flag.BoolVar(&staged, "staged", false, "Process only staged files (default behavior)")
	flag.BoolVar(&staged, "s", false, "Process only staged files (shorthand)")
	flag.BoolVar(&all, "all", false, "Process all files recursively (be careful with large codebases)")
	flag.BoolVar(&all, "a", false, "Process all files recursively (shorthand)")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.Parse()

	args := flag.Args()
	if len(args) > 0 && args[0] == "install-hooks" {
		err := InstallPreCommitHook(verbose)
		if err != nil {
			fmt.Printf("Error installing pre-commit hook: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Pre-commit hook installed successfully!")
		return
	}

	if showVersion {
		fmt.Printf("nocmt version %s\n", Version)
		return
	}

	commentConfig := config.New()
	err := commentConfig.LoadConfigurations()
	if err != nil && verbose {
		fmt.Printf("Warning: Could not load configuration: %v\n", err)
	}

	if configAdd != "" {
		err := commentConfig.AddIgnorePattern(configAdd)
		if err != nil {
			fmt.Printf("Error adding pattern to local config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Pattern '%s' added to local ignore list\n", configAdd)
		os.Exit(0)
	}

	if configAddGlobal != "" {
		err := commentConfig.AddGlobalIgnorePattern(configAddGlobal)
		if err != nil {
			fmt.Printf("Error adding pattern to global config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Pattern '%s' added to global ignore list\n", configAddGlobal)
		os.Exit(0)
	}

	if configAddFileIgnore != "" {
		err := commentConfig.AddFileIgnorePattern(configAddFileIgnore)
		if err != nil {
			fmt.Printf("Error adding pattern to local file ignore list: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Pattern '%s' added to local file ignore list\n", configAddFileIgnore)
		os.Exit(0)
	}

	if configAddFileIgnoreGlobal != "" {
		err := commentConfig.AddGlobalFileIgnorePattern(configAddFileIgnoreGlobal)
		if err != nil {
			fmt.Printf("Error adding pattern to global file ignore list: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Pattern '%s' added to global file ignore list\n", configAddFileIgnoreGlobal)
		os.Exit(0)
	}

	if ignorePatterns != "" {
		patterns := strings.Split(ignorePatterns, ",")
		for i := range patterns {
			patterns[i] = strings.TrimSpace(patterns[i])
		}
		err := commentConfig.SetCLIPatterns(patterns)
		if err != nil {
			fmt.Printf("Error parsing ignore patterns: %v\n", err)
			os.Exit(1)
		}
	}

	if ignoreFilePatterns != "" {
		patterns := strings.Split(ignoreFilePatterns, ",")
		for i := range patterns {
			patterns[i] = strings.TrimSpace(patterns[i])
		}
		err := commentConfig.SetCLIFilePatterns(patterns)
		if err != nil {
			fmt.Printf("Error parsing file ignore patterns: %v\n", err)
			os.Exit(1)
		}
	}

	inputPath := ""
	if len(args) > 0 {
		inputPath = args[0]
	}

	preserveDirectives = !removeDirectives

	if all {
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}

		if !walker.ValidateGitRepository(currentDir, force) {
			fmt.Println("Aborted: User chose not to proceed with non-git repository.")
			os.Exit(1)
		}

		processDirectory(currentDir, preserveDirectives, dryRun, verbose, force, commentConfig)
		return
	}

	if inputPath == "" {
		staged = true
	}

	if staged {
		if !isGitRepo() {
			fmt.Println("Error: can only process staged files inside a git repository")
			os.Exit(1)
		}
		processStagedFiles(preserveDirectives, dryRun, verbose, commentConfig)
		return
	}

	if inputPath != "" {
		fileInfo, err := os.Stat(inputPath)
		if os.IsNotExist(err) {
			fmt.Printf("Error: Path '%s' does not exist\n", inputPath)
			os.Exit(1)
		}

		if !fileInfo.IsDir() {
			processSingleFile(inputPath, preserveDirectives, commentConfig)
			return
		}

		if !walker.ValidateGitRepository(inputPath, force) {
			fmt.Println("Aborted: User chose not to proceed with non-git repository.")
			os.Exit(1)
		}

		processDirectory(inputPath, preserveDirectives, dryRun, verbose, force, commentConfig)
		return
	}

	fmt.Println("Error: No action specified")
	fmt.Println("Usage: nocmt [path] [options]")
	fmt.Println("       nocmt install-hooks")
	os.Exit(1)
}

func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil
}

func getStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only", "--diff-filter=ACM")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get staged files: %w", err)
	}

	files := strings.Split(strings.TrimSpace(string(out)), "\n")
	var result []string
	for _, file := range files {
		if file != "" {
			result = append(result, file)
		}
	}
	return result, nil
}

func getStagedFileContent(filePath string) (string, error) {
	cmd := exec.Command("git", "show", ":"+filePath)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get staged content for %s: %w", filePath, err)
	}
	return string(out), nil
}

func getModifiedLines(filePath string) (map[int]bool, error) {
	cmd := exec.Command("git", "diff", "--cached", "--unified=0", filePath)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get diff for %s: %w", filePath, err)
	}

	modifiedLines := make(map[int]bool)
	lineRegex := regexp.MustCompile(`^@@ -\d+(?:,\d+)? \+(\d+)(?:,(\d+))? @@`)

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		matches := lineRegex.FindStringSubmatch(line)
		if len(matches) >= 3 {
			lineNum, _ := strconv.Atoi(matches[1])
			count := 1
			if matches[2] != "" {
				count, _ = strconv.Atoi(matches[2])
			}
			for i := 0; i < count; i++ {
				modifiedLines[lineNum+i] = true
			}
		}
	}

	return modifiedLines, nil
}

func processStagedFiles(preserveDirectives bool, dryRun bool, verbose bool, commentConfig *config.Config) {
	factory := processor.NewProcessorFactory()
	factory.SetPreserveDirectives(preserveDirectives)
	factory.SetCommentConfig(commentConfig)

	stagedFiles, err := getStagedFiles()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if len(stagedFiles) == 0 {
		fmt.Println("No staged files found.")
		return
	}

	if verbose {
		fmt.Printf("Found %d staged files to process\n", len(stagedFiles))
	}

	processed := 0
	skipped := 0
	errors := 0

	for _, filePath := range stagedFiles {
		if verbose {
			fmt.Printf("Examining %s...\n", filePath)
		}

		proc, err := factory.GetProcessorByExtension(filePath)
		if err != nil {
			if verbose {
				fmt.Printf("Skipping %s: %v\n", filePath, err)
			}
			skipped++
			continue
		}

		stagedContent, err := getStagedFileContent(filePath)
		if err != nil {
			fmt.Printf("Error reading staged content of %s: %v\n", filePath, err)
			errors++
			continue
		}

		modifiedLines, err := getModifiedLines(filePath)
		if err != nil {
			fmt.Printf("Error getting modified lines for %s: %v\n", filePath, err)
			errors++
			continue
		}

		if verbose {
			fmt.Printf("Found %d modified lines in %s\n", len(modifiedLines), filePath)
		}

		if len(modifiedLines) == 0 {
			if verbose {
				fmt.Printf("No modified lines in %s, skipping\n", filePath)
			}
			skipped++
			continue
		}

		result, err := processFileWithSelectiveCommentRemoval(stagedContent, filePath, proc, modifiedLines, preserveDirectives, commentConfig)
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", filePath, err)
			if verbose {
				fmt.Printf("Note: Tree-sitter parsing is now required; no fallback option available\n")
			}
			errors++
			continue
		}

		if result == stagedContent {
			if verbose {
				fmt.Printf("No changes needed for %s\n", filePath)
			}
			skipped++
			continue
		}

		if verbose {
			fmt.Printf("Processing %s\n", filePath)
		}

		if !dryRun {
			err = os.WriteFile(filePath, []byte(result), 0644)
			if err != nil {
				fmt.Printf("Error writing file %s: %v\n", filePath, err)
				errors++
				continue
			}

			cmd := exec.Command("git", "add", filePath)
			err = cmd.Run()
			if err != nil {
				fmt.Printf("Error re-staging file %s: %v\n", filePath, err)
				errors++
				continue
			}

			if verbose {
				fmt.Printf("Successfully processed and re-staged %s\n", filePath)
			}
		} else if verbose {
			fmt.Printf("Dry run: would process %s\n", filePath)
		}

		processed++
	}

	fmt.Printf("\nProcessing complete:\n")
	fmt.Printf("- Files processed: %d\n", processed)
	fmt.Printf("- Files skipped: %d\n", skipped)
	fmt.Printf("- Errors: %d\n", errors)
}

func processFileWithSelectiveCommentRemoval(content string, filePath string, proc processor.LanguageProcessor, modifiedLines map[int]bool, preserveDirectives bool, commentConfig *config.Config) (string, error) {
	return processor.SelectivelyStripComments(content, filePath, proc, modifiedLines, preserveDirectives, commentConfig)
}

func processSingleFile(inputFile string, preserveDirectives bool, commentConfig *config.Config) {
	factory := processor.NewProcessorFactory()
	factory.SetPreserveDirectives(preserveDirectives)
	factory.SetCommentConfig(commentConfig)

	proc, err := factory.GetProcessorByExtension(inputFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	content, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	if proc.GetLanguageName() == "go" && preserveDirectives {
		proc = processor.NewGoProcessor(true)
		proc.SetCommentConfig(commentConfig)
	}

	if proc.GetLanguageName() == "bash" && preserveDirectives {
		proc = processor.NewBashProcessor(true)
		proc.SetCommentConfig(commentConfig)
	}

	result, err := proc.StripComments(string(content))
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(result)
}

func processDirectory(dirPath string, preserveDirectives bool, dryRun bool, verbose bool, force bool, commentConfig *config.Config) {
	config := walker.ProcessorConfig{
		PreserveDirectives: preserveDirectives,
		DryRun:             dryRun,
		Verbose:            verbose,
		Force:              force,
		CommentConfig:      commentConfig,
	}

	processorIntegration := walker.NewProcessorIntegration(config)

	fmt.Printf("Processing directory: %s\n", dirPath)
	if dryRun {
		fmt.Println("Running in dry-run mode - no changes will be written")
	}

	err := processorIntegration.ProcessRepository(dirPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	processed, skipped, errors := processorIntegration.GetStats()
	fmt.Printf("\nProcessing complete:\n")
	fmt.Printf("- Files processed: %d\n", processed)
	fmt.Printf("- Files skipped: %d\n", skipped)
	fmt.Printf("- Errors: %d\n", errors)
}
