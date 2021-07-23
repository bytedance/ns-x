package networksimulator

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	testTime := Now()
	assert.LessOrEqual(t, time.Since(testTime), time.Second)
}

func BenchmarkTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Now()
	}
}
