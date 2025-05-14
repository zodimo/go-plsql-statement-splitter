#!/bin/bash

# Make sure the script exits on any error
set -e

echo "=== Internal PL/SQL Parser Generator ==="


alias antlr4='java -Xmx500M -cp "./antlr-4.13.2-complete.jar:$CLASSPATH" org.antlr.v4.Tool'

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "ERROR: Go is not installed or not in PATH. Please install it first."
    exit 1
fi

echo "Go found: $(go version)"

# Verify the grammar file exists
GRAMMAR_FILE="PlSqlParser.g4"
if [ ! -f "$GRAMMAR_FILE" ]; then
    echo "ERROR: Grammar file not found: $GRAMMAR_FILE"
    exit 1
fi
echo "Grammar file found: $GRAMMAR_FILE"

# Set up output directory
mkdir -p gen
echo "Output directory: $(pwd)/gen"

# Generate parser
echo "Generating parser from PlSqlParser.g4..."
# antlr4 -Dlanguage=Go -package parser -o gen PlSqlParser.g4

antlr4 -Dlanguage=Go -package gen  -o gen -visitor *.g4

# Verify generation was successful
if [ $? -eq 0 ] && [ -d "gen" ] && [ "$(ls -A gen)" ]; then
    echo "Parser generated successfully!"
    echo "Generated files:"
    ls -la gen/
else
    echo "ERROR: Parser generation failed"
    exit 1
fi

echo "Parser generation completed successfully." 