# Product Requirements Document: go-plsql-splitter

## 1. Overview and Objectives

The go-plsql-splitter is a Go library designed to accurately split Oracle PL/SQL scripts into individual statements with precise boundary detection. Using ANTLR4 for parsing (with no regex allowed), it aims to provide developers with a reliable tool for extracting SQL statements that can be validated and executed in deployment pipelines.

### Key Objectives:
- Achieve 100% accurate PL/SQL statement boundary detection
- Provide precise source location tracking for each statement
- Deliver high performance for processing large SQL scripts
- Supply detailed error information for syntax issues

## 2. Target Audience

- Software developers building Oracle database deployment tools
- DevOps engineers implementing CI/CD pipelines for database changes
- Database administrators automating SQL script validation
- Anyone building tools that require precise PL/SQL statement extraction

## 3. Core Features and Functionality

### 3.1 Statement Splitting
- Split PL/SQL scripts into individual statements with 100% boundary accuracy
- Support all Oracle 19c PL/SQL statement types
- Preserve original statement content exactly as it appears in source

### 3.2 Position Tracking
- Track line and column numbers for the start and end of each statement
- Maintain context information for error reporting

### 3.3 Comment Handling
- Properly handle both single-line (--) and multi-line (/* */) comments
- Preserve comments in output when they're part of a statement

### 3.4 Input Processing
- Process input from file paths
- Process input from string content
- Consider streaming implementation for large files

### 3.5 Error Reporting
- Provide detailed syntax error messages
- Include file, line, and column information in error reports
- Offer context about the statement where the error occurred

### 3.6 Serialization Support
- All output structures must support JSON marshalling
- Consistent field naming for easy integration with other tools

## 4. Technical Requirements

### 4.1 Development
- Go 1.24 or later required
- ANTLR4 grammar for Oracle PL/SQL (supporting Oracle 19c syntax)
- No regex allowed for statement splitting
- Public GitHub repository under MIT license

### 4.2 ANTLR4 Integration
- May require creating or updating ANTLR grammar to support latest language features
- Custom grammar implementation if existing packages are outdated
- Performance optimization for the ANTLR4 parsing process

## 5. Performance Requirements

### 5.1 Processing Speed
- Optimize for parsing speed while maintaining accuracy
- Should handle multi-megabyte SQL scripts efficiently

### 5.2 Memory Efficiency
- Minimize memory usage, especially for large files
- Consider streaming approach for very large files to avoid loading entire content

### 5.3 Resource Utilization
- Avoid excessive CPU usage during parsing
- Balance accuracy with performance

## 6. Input/Output Specifications

### 6.1 Input
- File paths to PL/SQL scripts
- String content containing PL/SQL statements
- File encoding defaults to UTF-8

### 6.2 Output
- Structured data containing individual SQL statements
- Location information (line/column) for each statement
- Statement type classification if available from ANTLR parser
- Format that supports JSON marshalling

Example output structure:
```go
type Statement struct {
  Content     string `json:"content"`
  StartLine   int    `json:"startLine"`
  EndLine     int    `json:"endLine"`
  StartColumn int    `json:"startColumn"`
  EndColumn   int    `json:"endColumn"`
  Type        string `json:"type,omitempty"` // If available from ANTLR parser
}
```

## 7. Error Handling

### 7.1 Syntax Errors
- Detailed error messages for syntax issues
- Include file, line, and column information
- Clear description of the error nature

Example error structure:
```go
type SyntaxError struct {
  Message   string `json:"message"`
  Line      int    `json:"line"`
  Column    int    `json:"column"`
  Statement string `json:"statement,omitempty"`
}
```

### 7.2 File Errors
- Proper handling of file not found, permission issues, etc.
- Clear error messages for I/O problems

### 7.3 No Recovery Mechanism
- No error recovery mechanisms required
- Parser should fail with clear error when syntax is invalid

## 8. Integration with External Components

### 8.1 plsql-parser Integration
- Design for compatibility with github.com/zodimo/plsql-parser
- Clear interfaces for integration with this package

### 8.2 Deployment System Integration
- Output format suitable for passing to deployment systems
- Consider common deployment tool requirements

## 9. API Design

### 9.1 Simple Interface
```go
// Basic functions for simple use cases
func SplitFile(filePath string) ([]Statement, error)
func SplitString(content string) ([]Statement, error)
```

### 9.2 Configurable Interface
```go
// More flexible interface with configuration options
type Splitter struct {
  // Configuration options
}

func NewSplitter(options ...Option) *Splitter
func (s *Splitter) SplitFile(filePath string) ([]Statement, error)
func (s *Splitter) SplitString(content string) ([]Statement, error)
```

## 10. Testing and Quality Assurance

### 10.1 Test Coverage
- Comprehensive test suite with high coverage
- Unit tests for parser components
- Integration tests for full splitting functionality

### 10.2 Test Cases
- Tests for all Oracle 19c statement types
- Edge cases (comments, nested statements, etc.)
- Error cases with invalid syntax

### 10.3 Performance Testing
- Benchmark tests for performance optimization
- Memory usage monitoring
- Load testing with large SQL scripts

## 11. Deployment and Distribution

### 11.1 Packaging
- Standard Go module
- Public GitHub repository under user zodimo
- MIT license

### 11.2 Documentation
- Comprehensive README with usage examples
- Godoc API documentation
- Examples for common use cases

### 11.3 Versioning
- Semantic versioning (SemVer)
- Backward compatibility guarantees

## 12. Implementation Considerations

### 12.1 ANTLR4 Grammar
- Evaluate existing PL/SQL grammars for ANTLR4
- May need to fork and modify an existing grammar to support Oracle 19c
- Consider performance optimizations in the grammar

### 12.2 Parsing Strategy
- Use ANTLR4's parse tree listeners or visitors
- Track statement boundaries during parsing
- Handle special cases like anonymous blocks, stored procedures, etc.

### 12.3 Memory Management
- Avoid loading entire files into memory for large scripts
- Consider streaming parser implementation for large files
- Efficient string handling

## 13. Development Roadmap

### 13.1 Phase 1: Basic Functionality
- Set up project structure
- Implement basic file/string parsing
- Handle simple statement types

### 13.2 Phase 2: Enhanced Features
- Support all Oracle 19c statement types
- Implement detailed error reporting
- Add position tracking

### 13.3 Phase 3: Performance Optimization
- Benchmark and optimize for speed
- Memory usage optimization
- Handle edge cases

### 13.4 Phase 4: Documentation and Examples
- Comprehensive documentation
- Usage examples
- Integration examples

## 14. Challenges and Risks

### 14.1 ANTLR4 Grammar Complexity
- PL/SQL grammar is complex and may require significant effort to get 100% correct
- Specific Oracle 19c syntax features might be challenging to parse

### 14.2 Performance Optimization
- Balancing parsing accuracy with performance
- ANTLR4 parsers can be memory-intensive for complex grammars

### 14.3 Edge Cases
- Handling non-standard PL/SQL syntax extensions
- Correctly parsing complex nested constructs

## 15. Dependencies

### 15.1 ANTLR4 Go Runtime
- Required for parsing
- May have its own version constraints

### 15.2 plsql-parser Package
- Need to define clear interfaces with this package
- Ensure compatible design decisions

## 16. Future Considerations

### 16.1 Version Support
- Potential support for newer Oracle versions after 19c
- Backward compatibility with older Oracle versions

### 16.2 Feature Extensions
- Potential support for other SQL dialects
- Statement validation capabilities
- SQL transformation features

### 16.3 Performance Enhancements
- Ongoing optimization for even better performance
- Support for concurrent parsing 