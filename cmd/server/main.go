package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"yuxialuozi_graduation_design_backend/internal/wire"
)

func main() {
	// 初始化日志
	initLogger()

	// 初始化应用
	router, cleanup, err := wire.InitializeApp()
	if err != nil {
		zap.L().Fatal("Failed to initialize app", zap.Error(err))
	}
	defer cleanup()

	zap.L().Info("Server starting on :8080")

	if err := router.Run(); err != nil {
		zap.L().Fatal("Failed to start server", zap.Error(err))
	}
}

func initLogger() {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
}
