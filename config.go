package ns_x

import (
	"time"
)

// DefaultBucketSize used for simulation
const DefaultBucketSize = time.Second

// DefaultMaxBuckets used for simulation
const DefaultMaxBuckets = 128

// Config something of the simulation, usually for optimization
type Config func(config *config)

type config struct {
	bucketSize time.Duration
	maxBuckets int
}

// WithBucketSize set the bucket size of each bucket, usually used with WithMaxBuckets
// event queue use a bucket sort to separate events into buckets, and then use a heap sort for each bucket
// maxBuckets * bucketSize should cover events cluster, usually related to the simulation content
// for example, most events generated with 1-second delay, then maxBuckets * bucketSize should be greater than 1 second
// count of events in each bucket should not be too large
func WithBucketSize(bucketSize time.Duration) Config {
	return func(config *config) {
		config.bucketSize = bucketSize
	}
}

// WithMaxBuckets set the max bucket count, for detail see WithBucketSize
func WithMaxBuckets(maxBuckets int) Config {
	return func(config *config) {
		config.maxBuckets = maxBuckets
	}
}

func (c *config) apply(configs ...Config) {
	for _, config := range configs {
		config(c)
	}
}
