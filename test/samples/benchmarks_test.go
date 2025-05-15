package samples

import (
	"testing"
)

// BenchmarkGetSimpleSQLSamples benchmarks the GetSimpleSQLSamples function
func BenchmarkGetSimpleSQLSamples(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetSimpleSQLSamples()
	}
}

// BenchmarkGetComplexSQLSamples benchmarks the GetComplexSQLSamples function
func BenchmarkGetComplexSQLSamples(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetComplexSQLSamples()
	}
}

// BenchmarkGetInvalidSQLSamples benchmarks the GetInvalidSQLSamples function
func BenchmarkGetInvalidSQLSamples(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetInvalidSQLSamples()
	}
}
