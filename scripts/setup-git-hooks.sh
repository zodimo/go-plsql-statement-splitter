#!/bin/bash

# Creates git hooks for the project

HOOK_DIR=.git/hooks
PRE_COMMIT=${HOOK_DIR}/pre-commit

# Create hooks directory if it doesn't exist
mkdir -p ${HOOK_DIR}

# Create pre-commit hook
cat > ${PRE_COMMIT} << 'EOF'
#!/bin/bash

# Pre-commit hook for go-plsql-statement-splitter
# This hook runs tests before allowing a commit

echo "Running pre-commit tests..."

# Stash any changes not being committed
git stash -q --keep-index

# Run tests
./scripts/test.sh

# Store the result
RESULT=$?

# Restore stashed changes
git stash pop -q

# Return the test result
if [ $RESULT -ne 0 ]; then
  echo "Tests failed. Commit aborted."
  exit 1
fi

exit 0
EOF

# Make the pre-commit hook executable
chmod +x ${PRE_COMMIT}

echo "Git hooks installed successfully. Pre-commit hook will run tests before each commit." 