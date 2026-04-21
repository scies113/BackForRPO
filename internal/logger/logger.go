package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// Log - глобальный логгер приложения
var Log *zap.Logger

// Init - инициализация логгера
// В режиме development выводит читаемые логи, в production - JSON
func Init() {
	env := os.Getenv("APP_ENV")

	var cfg zap.Config

	if env == "production" {
		// Продакшн: JSON, уровень info
		cfg = zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	} else {
		// Разработка: читаемые логи, уровень debug
		cfg = zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var err error
	Log, err = cfg.Build()
	if err != nil {
		panic("Ошибка инициализации логгера: " + err.Error())
	}
}

// Info - логирование информационного сообщения
func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

// Error - логирование ошибки
func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

// Debug - логирование отладочного сообщения
func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

// Warn - логирование предупреждения
func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

// Fatal - логирование критической ошибки и завершение программы
func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}

// Sync - сброс буферов логгера (вызывать при завершении)
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
