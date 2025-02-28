package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

type Logger struct {
	log *zap.Logger
}

func GetLogger(requestID string, fields ...zapcore.Field) Logger {
	if logger == nil {
		setup()
	}

	l := Logger{log: logger.With(fields...)}
	if requestID != "" {
		l.log = logger.With(zap.String("X-Request-Id", requestID))
	}

	return l
}

func (l *Logger) SetFields(fields ...zapcore.Field) {
	l.log = l.log.With(fields...)
}

func (l *Logger) SetOptions(options ...zap.Option) {
	l.log = l.log.WithOptions(options...)
}

func (l *Logger) Error(msg string, fields ...zapcore.Field) {
	l.log.Error(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zapcore.Field) {
	l.log.Info(msg, fields...)
}

func setup() {
	allPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.DebugLevel
	})

	prodEncoder := zap.NewProductionEncoderConfig()
	prodEncoder.EncodeTime = zapcore.RFC3339TimeEncoder

	allCore := zapcore.NewCore(zapcore.NewJSONEncoder(prodEncoder), zapcore.Lock(os.Stdout), allPriority)

	logger = zap.New(zapcore.NewTee(allCore), zap.AddCaller(), zap.AddCallerSkip(3))
}
