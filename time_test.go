package networksimulator

import "testing"

func BenchmarkTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Now()
	}
}
