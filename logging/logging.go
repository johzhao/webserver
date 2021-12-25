package logging

import (
	"context"
	"fmt"
	"os"
	"webserver/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const TracerIDKey = "tracer-id"

func ContextField(ctx context.Context) []zap.Field {
	valueRaw := ctx.Value(TracerIDKey)
	value, ok := valueRaw.(string)
	if !ok {
		value = ""
	}

	return []zap.Field{
		{Key: TracerIDKey, Type: zapcore.StringType, String: value},
	}
}

func SetupLogger(config config.Logger) (*zap.Logger, error) {
	logLevel := zapcore.InfoLevel
	if err := logLevel.Set(config.Level); err != nil {
		return nil, fmt.Errorf("set log level (%s) failed with error: (%w)", config.Level, err)
	}

	atomLv := zap.NewAtomicLevel()
	atomLv.SetLevel(logLevel)

	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonEncoder := zapcore.NewJSONEncoder(cfg)

	consoleWriter := zapcore.AddSync(os.Stdout)

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.Filepath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxAge,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	})

	consoleCore := zapcore.NewCore(
		jsonEncoder,
		consoleWriter,
		atomLv,
	)

	fileCore := zapcore.NewCore(
		jsonEncoder,
		fileWriter,
		atomLv,
	)

	loggerCore := zapcore.NewTee(consoleCore, fileCore)

	return zap.New(loggerCore), nil
}
