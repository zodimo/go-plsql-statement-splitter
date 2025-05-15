# Testing Workflow

This document outlines the testing workflow for the go-plsql-statement-splitter project.

## Local Testing

### Running Tests

The project includes a test script that can be used to run tests locally. The script is located at `scripts/test.sh` and supports the following options:

```bash
./scripts/test.sh [flags]

Flags:
  -v              Verbose mode
  -r              Run tests with race detection
  -c              Generate coverage report
  -a              Run all checks (equivalent to -v -r -c)
```

The script will exit with a non-zero status code if any of the tests fail, making it suitable for use in CI/CD pipelines and Git hooks.

### Examples

Run basic tests:
```bash
./scripts/test.sh
```

Run tests in verbose mode with race detection:
```bash
./scripts/test.sh -v -r
```

Run all checks (tests, race detection, coverage):
```bash
./scripts/test.sh -a
```

### Coverage Reports

To generate a coverage report:
```bash
./scripts/test.sh -c
```

To view the coverage report in your browser:
```bash
go tool cover -html=coverage.out
```

## Git Hooks

The project uses Git hooks to ensure tests pass before commits are made. To set up Git hooks, run:

```bash
./scripts/setup-git-hooks.sh
```

This will install a pre-commit hook that runs tests before each commit. If tests fail, the commit will be aborted.

## Continuous Integration

The project uses GitHub Actions for continuous integration. The following workflows are configured:

1. **Test Workflow** - Runs tests on push to main branch and on pull requests
   - Runs standard tests
   - Runs tests with race detection
   - Performs linting

2. **Coverage Workflow** - Generates code coverage reports
   - Uploads coverage reports to Codecov

3. **Security Scanning** - Performs security checks
   - Runs gosec to detect security issues
   - Scans dependencies for vulnerabilities

## Adding Tests

When adding new features or fixing bugs, please include tests for the new functionality. Tests should be placed in the appropriate package directory with a `_test.go` suffix.

For example:
- Core package tests go in `pkg/splitter/*_test.go`
- Internal tests go in `internal/parser/*_test.go`
- Test samples go in `test/samples/`

## Test Sample Files

SQL sample files for testing are maintained in the `test/samples` directory. If you're adding support for new SQL constructs, please add sample files to test the new functionality.

## Benchmarks

Performance benchmarks are included in the test suite. To run benchmarks:

```bash
go test -bench=. ./...
```

To compare benchmark results before and after changes, you can use the benchcmp tool:

```bash
# Install benchcmp
go install golang.org/x/tools/cmd/benchcmp@latest

# Run benchmarks before changes and save results
go test -bench=. ./... > bench-before.txt

# Make changes, then run benchmarks again
go test -bench=. ./... > bench-after.txt

# Compare results
benchcmp bench-before.txt bench-after.txt
``` 