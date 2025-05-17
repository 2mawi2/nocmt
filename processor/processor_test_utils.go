package processor

import (
	"os"
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/assert"
)

func RunFileBasedTestCase(t *testing.T, processor LanguageProcessor, originalFilePath string, expectedFilePath string) {
	t.Helper()

	originalContentBytes, err := os.ReadFile(originalFilePath)
	if err != nil {
		t.Fatalf("Failed to read original file %s: %v", originalFilePath, err)
	}
	originalContent := string(originalContentBytes)

	expectedContentBytes, err := os.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("Failed to read expected file %s: %v", expectedFilePath, err)
	}
	expectedContent := string(expectedContentBytes)

	actualContent, err := processor.StripComments(originalContent)
	assert.NoError(t, err, "Processor StripComments failed for %s", originalFilePath)

	if expectedContent != actualContent {
		diff := difflib.UnifiedDiff{
			A:        difflib.SplitLines(expectedContent),
			B:        difflib.SplitLines(actualContent),
			FromFile: "Expected: " + expectedFilePath,
			ToFile:   "Actual",
			Context:  3,
		}
		diffStr, err := difflib.GetUnifiedDiffString(diff)
		if err != nil {
			t.Fatalf("Failed to generate diff: %v", err)
		}
		t.Errorf("Processed content does not match expected content for %s.\nDiff:\n%s", originalFilePath, diffStr)
	}
}

func normalizeNewlinesAndTrim(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	trimmedContent := strings.Join(lines, "\n")
	trimmedContent = strings.Trim(trimmedContent, "\n")
	if trimmedContent == "" {
		return ""
	}
	return trimmedContent + "\n"
}

func RunFileBasedTestCaseNormalized(t *testing.T, processor LanguageProcessor, originalFilePath string, expectedFilePath string) {
	t.Helper()

	originalContentBytes, err := os.ReadFile(originalFilePath)
	if err != nil {
		t.Fatalf("Failed to read original file %s: %v", originalFilePath, err)
	}
	originalContent := string(originalContentBytes)

	expectedContentBytes, err := os.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("Failed to read expected file %s: %v", expectedFilePath, err)
	}
	expectedContent := string(expectedContentBytes)

	actualContent, err := processor.StripComments(originalContent)
	assert.NoError(t, err, "Processor StripComments failed for %s", originalFilePath)

	normalizedExpected := normalizeNewlinesAndTrim(expectedContent)
	normalizedActual := normalizeNewlinesAndTrim(actualContent)

	if normalizedExpected != normalizedActual {
		diff := difflib.UnifiedDiff{
			A:        difflib.SplitLines(normalizedExpected),
			B:        difflib.SplitLines(normalizedActual),
			FromFile: "Expected (Normalized): " + expectedFilePath,
			ToFile:   "Actual (Normalized)",
			Context:  3,
		}
		diffStr, err := difflib.GetUnifiedDiffString(diff)
		if err != nil {
			t.Fatalf("Failed to generate diff: %v", err)
		}
		t.Errorf("Processed content (normalized) does not match expected content for %s.\nDiff:\n%s", originalFilePath, diffStr)
	}
}

func RunFileBasedTestCaseVeryLenient(t *testing.T, processor LanguageProcessor, originalFilePath string, expectedFilePath string) {
	t.Helper()

	originalContentBytes, err := os.ReadFile(originalFilePath)
	if err != nil {
		t.Fatalf("Failed to read original file %s: %v", originalFilePath, err)
	}
	originalContent := string(originalContentBytes)

	expectedContentBytes, err := os.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("Failed to read expected file %s: %v", expectedFilePath, err)
	}
	expectedContent := string(expectedContentBytes)

	actualContent, err := processor.StripComments(originalContent)
	assert.NoError(t, err, "Processor StripComments failed for %s", originalFilePath)

	normalizeFunc := func(s string) string {
		lines := strings.Split(s, "\n")
		var nonEmptyLines []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				nonEmptyLines = append(nonEmptyLines, trimmed)
			}
		}
		return strings.Join(nonEmptyLines, "\n")
	}

	normalizedExpected := normalizeFunc(expectedContent)
	normalizedActual := normalizeFunc(actualContent)

	if normalizedExpected != normalizedActual {
		t.Errorf("Processed content does not match expected content for %s after lenient normalization.", originalFilePath)
	}
}
