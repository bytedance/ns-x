package ns_x

import (
	"github.com/bytedance/ns-x/v2/base"
	"time"
)

type Config struct {
	BucketSize    time.Duration
	MaxBuckets    int
	InitialEvents []base.Event
}
