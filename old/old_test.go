package old

import (
	"testing"
)

func BenchmarkGetTransformDocuments(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetTransformDocuments()
	}
}