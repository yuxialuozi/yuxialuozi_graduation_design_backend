package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	_ "yuxialuozi_graduation_design_backend/docs"
	"yuxialuozi_graduation_design_backend/internal/wire"
)

// @title 租户信息管理系统 API
// @version 1.0
// @description 租户信息管理系统后端 API 服务
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 输入 Bearer {token}

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
