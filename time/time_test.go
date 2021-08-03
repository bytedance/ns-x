package time

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	for i := 0; i < 1000; i++ {
		testTime := Now()
		assert.LessOrEqual(t, time.Since(testTime), time.Microsecond*50)
	}
}

func BenchmarkTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Now()
	}
}
