# Contributing to nocmt

Thank you for considering contributing to nocmt! Here's how you can help.

## Development Setup

1. Clone the repository
   ```bash
   git clone https://github.com/2mawi2/nocmt.git
   cd nocmt
   ```

2. Install dependencies
   ```bash
   go mod download
   ```

3. Build the project
   ```bash
   go build
   ```

## Running Tests

Run tests using the test command:
```bash
just test
```

Or directly with Go:
```bash
go test ./...
```

## Running Benchmarks

```bash
just bench quick     # Fast benchmark (default)
just bench run       # Full benchmark suite (more accurate)
just bench compare   # Compare with baseline
just bench history   # Show benchmark history
```

## Adding Support for New Languages

To add support for a new language:

1. Create a new processor in `processor/` that implements the `LanguageProcessor` interface
2. Add the appropriate tree-sitter parser for the language
3. Register the processor in the `ProcessorFactory`
4. Add file extension mappings in the `GetProcessorByExtension` method

Example implementation for a new language:

```go
// MyLanguageProcessor implements LanguageProcessor for MyLanguage
type MyLanguageProcessor struct {
    preserveDirectives bool
}

// Implement all required methods...

// Register in factory
factory.Register(NewMyLanguageProcessor(false))

// Add extension mapping
// In GetProcessorByExtension:
extMap := map[string]string{
    // ...
    ".mylang": "mylanguage",
    // ...
}
```

## Pull Request Process

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Code Style

- Follow standard Go conventions and best practices
- Run `go fmt` before committing
- Add tests for new functionality

## Reporting Bugs

Please use GitHub Issues to report bugs, including:
- Clear description of the issue
- Steps to reproduce
- Expected vs. actual behavior
- Version information 