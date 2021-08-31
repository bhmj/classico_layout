package log

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ContextKey string

const (
	// ContextValueRequestID defines request ID value in context
	ContextValueRequestID ContextKey = "requestID"
	// ContextValueKeyID defines key ID value in context
	ContextValueKeyID ContextKey = "keyID"
)

// Logger implements logging functionality
type Logger interface {
	L() *zap.SugaredLogger
	LG() *zap.Logger
}

type logger struct {
	l  *zap.SugaredLogger
	lg *zap.Logger
}

// New returns new logger
func New(level string) (Logger, error) {
	// check level values
	zapLevels := map[string]zapcore.Level{
		"debug":  zap.DebugLevel,
		"info":   zap.InfoLevel,
		"warn":   zap.WarnLevel,
		"error":  zap.ErrorLevel,
		"dpanic": zap.DPanicLevel,
		"panic":  zap.PanicLevel,
		"fatal":  zap.FatalLevel,
	}

	var config zap.Config
	if level == "debug" {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
		config.DisableCaller = true
	}

	config.Level.SetLevel(zapLevels[level])
	config.Encoding = "json"
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.EncodeTime = timeEncoder

	lg, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("construct logger from config: %w", err)
	}
	return &logger{lg: lg, l: lg.Sugar()}, nil
}

func (l *logger) L() *zap.SugaredLogger {
	return l.l
}

func (l *logger) LG() *zap.Logger {
	return l.lg
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	encoded := t.UTC().AppendFormat([]byte{}, "2006-01-02T15:04:05.000Z")
	enc.AppendByteString(encoded)
}
