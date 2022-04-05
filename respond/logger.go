package respond

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func init() {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05.000")
	cfg.Encoding = "console"
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.DisableCaller = true
	cfg.DisableStacktrace = true

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = l.Sync()
	}()

	logger = l.Sugar()
}

// SetLogger sets the global Zap logger.
func SetLogger(l *zap.SugaredLogger) {
	logger = l
}
