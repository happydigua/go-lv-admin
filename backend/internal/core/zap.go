package core

import (
	"fmt"
	"os"

	"go-lv-vue-admin/internal/global"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Zap requires global.LV_CONFIG to be initialized
func Zap() (logger *zap.Logger) {
	if _, err := os.Stat(global.LV_CONFIG.Zap.Director); os.IsNotExist(err) {
		_ = os.Mkdir(global.LV_CONFIG.Zap.Director, os.ModePerm)
	}

	cores := internalZap()
	logger = zap.New(zapcore.NewTee(cores...))

	if global.LV_CONFIG.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}

func internalZap() []zapcore.Core {
	var cores []zapcore.Core

	fileEncoder := getEncoder()
	writer := getLogWriter()
	level := getLevel()

	fileCore := zapcore.NewCore(fileEncoder, writer, level)
	cores = append(cores, fileCore)

	if global.LV_CONFIG.Zap.LogInConsole {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		consoleCore := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)
		cores = append(cores, consoleCore)
	}

	return cores
}

func getEncoder() zapcore.Encoder {
	if global.LV_CONFIG.Zap.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  global.LV_CONFIG.Zap.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func getLogWriter() zapcore.WriteSyncer {
	file, _ := os.Create(fmt.Sprintf("%s/log.log", global.LV_CONFIG.Zap.Director))
	return zapcore.AddSync(file)
}

func getLevel() zapcore.Level {
	switch global.LV_CONFIG.Zap.Level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
