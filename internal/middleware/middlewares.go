package middleware

import (
	"crypto/md5"
	"fmt"
	"github.com/anilozgok/rate-limiter/internal/config"
	"github.com/anilozgok/rate-limiter/internal/limiter"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"math/big"
)

func RateLimiter(c *fiber.Ctx) error {
	globalRateLimiterConfig, err := config.Get()
	if err != nil {
		zap.L().Panic("failed to read configs", zap.Error(err))
	}

	tokenBucket := limiter.GetBucket(globalRateLimiterConfig, getIdentifier(c))
	if !tokenBucket.IsAllowed(1) {
		return c.SendStatus(fiber.StatusTooManyRequests)
	}
	return c.Next()
}

func getIdentifier(c *fiber.Ctx) string {
	ip := c.IP()
	url := c.OriginalURL()
	identifier := fmt.Sprintf("%s.%s", ip, url)
	h := md5.Sum([]byte(identifier))
	return new(big.Int).SetBytes(h[:]).Text(62)
}
