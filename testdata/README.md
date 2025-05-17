# nocmt Test Fixtures

This directory contains test fixtures for nocmt's language processors.

## Directory Structure

Each language has its own subdirectory under `testdata/` containing at minimum:

- `original.<ext>` - Source code with comments to be processed
- `expected.<ext>` - Expected output after comment removal

## File-Based Testing Pattern

The nocmt project uses file-based tests instead of inline table-driven tests for better maintainability:

1. The file-based test reads the `original.<ext>` file
2. It processes it using the appropriate language processor
3. It compares the result against `expected.<ext>`

## Adding New Test Cases

When adding support for a new language or new comment edge cases:

1. Create a new directory `testdata/your_language/` if it doesn't exist
2. Add an `original.<ext>` file with representative examples of all comment types
3. Generate or manually create the expected output in `expected.<ext>`
4. Add a test in `your_language_processor_test.go` that uses `RunFileBasedTestCaseNormalized`

## Special Cases

Specialized cases that require additional testing (like directive handling, errors, etc.)
should still be tested with traditional in-code test cases in the processor test files.
The file-based tests are primarily for testing comment-stripping capabilities. 