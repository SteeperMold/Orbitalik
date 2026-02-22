package zaplog

import (
	"time"

	"github.com/SteeperMold/Orbitalik/common/go/log"
	"go.uber.org/zap"
)

type Logger struct {
	l *zap.Logger
}

func New(l *zap.Logger) *Logger {
	return &Logger{l: l}
}

func (z *Logger) Info(msg string, fields ...log.Field) {
	z.l.Info(msg, toZap(fields)...)
}

func (z *Logger) Error(msg string, fields ...log.Field) {
	z.l.Error(msg, toZap(fields)...)
}

func toZap(fields []log.Field) []zap.Field {
	out := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		switch v := f.Value.(type) {
		case string:
			out = append(out, zap.String(f.Key, v))
		case int:
			out = append(out, zap.Int(f.Key, v))
		case int64:
			out = append(out, zap.Int64(f.Key, v))
		case bool:
			out = append(out, zap.Bool(f.Key, v))
		case float64:
			out = append(out, zap.Float64(f.Key, v))
		case time.Duration:
			out = append(out, zap.Duration(f.Key, v))
		default:
			out = append(out, zap.Any(f.Key, v))
		}
	}
	return out
}
