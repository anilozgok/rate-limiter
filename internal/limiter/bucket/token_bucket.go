package bucket

import (
	"math"
	"sync"
	"time"
)

type TokenBucket struct {
	rate           int64
	maxTokens      int64
	currentTokens  int64
	lastRefillTime time.Time
	mutex          sync.Mutex
}

func NewTokenBucket(rate int64, maxTokens int64) *TokenBucket {
	return &TokenBucket{
		rate:           rate,
		maxTokens:      maxTokens,
		currentTokens:  maxTokens,
		lastRefillTime: time.Now(),
	}
}

func (tb *TokenBucket) refill() {
	now := time.Now()
	end := time.Since(tb.lastRefillTime)
	tokensToAdd := tb.rate * int64(end.Seconds())
	tb.currentTokens = int64(math.Min(float64(tb.currentTokens+tokensToAdd), float64(tb.maxTokens)))
	tb.lastRefillTime = now
}

func (tb *TokenBucket) IsAllowed(count int64) bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	tb.refill()
	if tb.currentTokens >= count {
		tb.currentTokens -= count
		return true
	}
	return false
}
