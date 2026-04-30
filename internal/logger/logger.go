package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func NewLogger(environment string) *Logger {
	var config zap.Config

	if environment == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	log, _ := config.Build()

	return &Logger{log}
}

func (l *Logger) Info(msg string, fields map[string]interface{}) {
	zapFields := fieldsToZapFields(fields)
	l.Logger.Info(msg, zapFields...)
}

func (l *Logger) Error(msg string, fields map[string]interface{}) {
	zapFields := fieldsToZapFields(fields)
	l.Logger.Error(msg, zapFields...)
}

func (l *Logger) Warn(msg string, fields map[string]interface{}) {
	zapFields := fieldsToZapFields(fields)
	l.Logger.Warn(msg, zapFields...)
}

func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	zapFields := fieldsToZapFields(fields)
	l.Logger.Debug(msg, zapFields...)
}

func (l *Logger) Fatal(msg string, fields map[string]interface{}) {
	zapFields := fieldsToZapFields(fields)
	l.Logger.Fatal(msg, zapFields...)
}

func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

func fieldsToZapFields(fields map[string]interface{}) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}
