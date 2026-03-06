package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init(level string, outputPaths []string) error {
	// 确保每个输出路径的目录存在
	for _, path := range outputPaths {
		if path == "stdout" || path == "stderr" {
			continue
		}
		dir := filepath.Dir(path)
		if dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("创建日志目录失败 %s: %w", dir, err)
			}
		}
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(getLevel(level))
	cfg.OutputPaths = outputPaths
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	l, err := cfg.Build()
	if err != nil {
		return err
	}
	log = l
	return nil
}

func getLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// L 提供全局 logger 访问
func L() *zap.Logger {
	if log == nil {
		// 如果没有初始化，则创建一个默认的（仅用于防止 panic）
		l, _ := zap.NewProduction()
		return l
	}
	return log
}

// S 如果需要 SugarLogger
func S() *zap.SugaredLogger {
	return L().Sugar()
}
