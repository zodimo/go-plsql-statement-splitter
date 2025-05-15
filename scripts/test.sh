#!/bin/bash

# This script runs all tests locally for the go-plsql-statement-splitter project
# Usage: ./scripts/test.sh [flags]
#
# Flags:
#   -v              Verbose mode
#   -r              Run tests with race detection
#   -c              Generate coverage report
#   -l              Run linter (currently disabled)
#   -a              Run all checks (equivalent to -v -r -c -l)

VERBOSE=0
RACE=0
COVERAGE=0
LINT=0  # Linting is disabled until further notice
EXIT_CODE=0

while getopts "vrcla" opt; do
  case $opt in
    v) VERBOSE=1 ;;
    r) RACE=1 ;;
    c) COVERAGE=1 ;;
    l) echo "Note: Linting is currently disabled until further notice." ;;
    a) VERBOSE=1; RACE=1; COVERAGE=1; ;;  # LINT=1 removed
    *) echo "Invalid option: -$OPTARG" >&2; exit 1 ;;
  esac
done

# Function to run test command and capture exit code
run_test() {
  echo "Running:" "$@"
  "$@"
  local result=$?
  if [ $result -ne 0 ]; then
    EXIT_CODE=$result
    echo "FAILED: Command exited with code $result" >&2
  fi
  return $result
}

# Determine test flags
TEST_FLAGS=""
if [ $VERBOSE -eq 1 ]; then
  TEST_FLAGS="-v"
fi

# LINTING FUNCTIONALITY DISABLED UNTIL FURTHER NOTICE
# Run linter if requested
#if [ $LINT -eq 1 ]; then
#  echo "=== Running golangci-lint ==="
#  
#  # Check if golangci-lint is installed
#  if ! command -v golangci-lint &> /dev/null; then
#    echo "golangci-lint not found, attempting to install..."
#    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(go env GOPATH)"/bin v1.55.0
#    
#    # Check if the installation was successful
#    if ! command -v "$(go env GOPATH)"/bin/golangci-lint &> /dev/null; then
#      echo "ERROR: Failed to install golangci-lint" >&2
#      EXIT_CODE=1
#    fi
#  fi
#  
#  # Create .golangci.yml file if it doesn't exist
#  if [ ! -f .golangci.yml ]; then
#    echo "Creating golangci-lint config..."
#    cat > .golangci.yml << EOL
#run:
#  timeout: 5m
#  tests: true
#  skip-dirs:
#    - internal/parser/gen
#  
#linters:
#  disable-all: true
#  enable:
#    - errcheck
#    - gosimple
#    - govet
#    - ineffassign
#    - staticcheck
#    - unused
#EOL
#  fi
#  
#  # Run golangci-lint
#  echo "Running linter..."
#  if command -v golangci-lint &> /dev/null; then
#    run_test golangci-lint run -v
#  else
#    run_test "$(go env GOPATH)"/bin/golangci-lint run -v
#  fi
#fi

echo "=== Running Go tests ==="
if [ $RACE -eq 1 ]; then
  echo "Running tests with race detection..."
  if [ -n "$TEST_FLAGS" ]; then
    run_test go test "$TEST_FLAGS" -race ./...
  else
    run_test go test -race ./...
  fi
else
  if [ -n "$TEST_FLAGS" ]; then
    run_test go test "$TEST_FLAGS" ./...
  else
    run_test go test ./...
  fi
fi

if [ $COVERAGE -eq 1 ]; then
  echo "=== Generating coverage report ==="
  run_test go test -coverprofile=coverage.out -covermode=atomic ./...
  
  # Only run cover tool if the previous command succeeded
  if [ $? -eq 0 ]; then
    go tool cover -func=coverage.out
    echo "For HTML coverage report, run: go tool cover -html=coverage.out"
  fi
fi

echo "=== All tests completed ==="

# Return the captured exit code
if [ $EXIT_CODE -ne 0 ]; then
  echo "ERROR: One or more tests failed. Please check the output above." >&2
fi

exit $EXIT_CODE 