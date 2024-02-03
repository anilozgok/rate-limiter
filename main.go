package main

import (
	"fmt"
	"github.com/anilozgok/rate-limiter/internal/config"
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
	_, err := config.Get()
	if err != nil {
		zap.L().Fatal("failed to read configs", zap.Error(err))
	}

}
