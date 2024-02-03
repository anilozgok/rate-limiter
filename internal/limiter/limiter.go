package limiter

import (
	"github.com/anilozgok/rate-limiter/internal/config"
	"github.com/anilozgok/rate-limiter/internal/limiter/bucket"
)

var clientBucketMap = make(map[string]*bucket.TokenBucket)

func GetBucket(config *config.GlobalRateLimiter, identifier string) *bucket.TokenBucket {
	if clientBucketMap[identifier] == nil {
		clientBucketMap[identifier] = bucket.NewTokenBucket(int64(config.RateLimiter.Limits.Rate), int64(config.RateLimiter.Limits.Count))
	}
	return clientBucketMap[identifier]
}
