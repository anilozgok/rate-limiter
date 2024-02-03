package main

import (
	"fmt"
	"github.com/anilozgok/rate-limiter/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func init() {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	l, err := cfg.Build()
	if err != nil {
		panic(fmt.Sprintf("fail to build log. err: %s", err))
	}

	zap.ReplaceGlobals(l)
}

func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(logger.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Get("/non-limited-endpoint", func(c *fiber.Ctx) error {
		return c.JSON("non-limited-endpoint")
	})

	app.Get("/limited-endpoint", middleware.RateLimiter, func(c *fiber.Ctx) error {
		return c.JSON("limited-endpoint")
	})

	zap.L().Info("server starting on :8080")
	if err := app.Listen(":8080"); err != nil {
		zap.L().Panic("error occurred while starting server", zap.Error(err))
	}

	if err := app.Shutdown(); err != nil {
		zap.L().Panic("error occurred while shutting down server", zap.Error(err))
	}
}
