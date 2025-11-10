package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wrapper para zap.Logger
type Logger struct {
	*zap.Logger
}

// New crea una nueva instancia del logger
func New(environment string) (*Logger, error) {
	var config zap.Config

	if environment == "production" {
		// Producción: JSON estructurado
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		// Development: legible para humanos
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Nivel de log
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if os.Getenv("LOG_LEVEL") == "debug" {
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	// Salida
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build(
		zap.AddCallerSkip(1), // Skip wrapper
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, err
	}

	return &Logger{Logger: logger}, nil
}

// WithFields agrega campos contextuales
func (l *Logger) WithFields(fields ...zap.Field) *Logger {
	return &Logger{Logger: l.Logger.With(fields...)}
}

// WithRequestID agrega el request ID al logger
func (l *Logger) WithRequestID(requestID string) *Logger {
	return l.WithFields(zap.String("request_id", requestID))
}

// WithUserID agrega el user ID al logger
func (l *Logger) WithUserID(userID int64) *Logger {
	return l.WithFields(zap.Int64("user_id", userID))
}

// InfoWithFields log info con campos
func (l *Logger) InfoWithFields(msg string, fields ...zap.Field) {
	l.Info(msg, fields...)
}

// ErrorWithFields log error con campos
func (l *Logger) ErrorWithFields(msg string, fields ...zap.Field) {
	l.Error(msg, fields...)
}

// WarnWithFields log warning con campos
func (l *Logger) WarnWithFields(msg string, fields ...zap.Field) {
	l.Warn(msg, fields...)
}

// DebugWithFields log debug con campos
func (l *Logger) DebugWithFields(msg string, fields ...zap.Field) {
	l.Debug(msg, fields...)
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}
