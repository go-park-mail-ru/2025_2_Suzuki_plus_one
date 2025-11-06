// logger/zap_logger.go
package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Ensure zapLogger implements Logger
var _ Logger = (*zapLogger)(nil)

type zapLogger struct {
	logger *zap.Logger
}

// Create a new Logger instance using zap
//
// Panics if zap logger initialization fails
func NewZapLogger(development bool) Logger {
	var zapLog *zap.Logger
	var err error

	cfg := zap.NewDevelopmentConfig()
	if !development {
		cfg = zap.NewProductionConfig()
	}
	zapLog, err = cfg.Build(
		zap.AddCaller(), zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.PanicLevel), // only for Error and above
	)

	if err != nil {
		panic("Can't initialize zap logger: " + err.Error())
	}

	return &zapLogger{logger: zapLog}
}

// OPTIMIZATION: The wrapper method convertToZapFields can slow down the performance

func (z *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	fields := convertToZapFields(keysAndValues...)
	z.logger.Debug(msg, fields...)
}

// Zap Info wrapper
func (z *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	fields := convertToZapFields(keysAndValues...)
	z.logger.Info(msg, fields...)
}

// Zap Warn wrapper
func (z *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	fields := convertToZapFields(keysAndValues...)
	z.logger.Warn(msg, fields...)
}

// Zap Error wrapper
func (z *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	fields := convertToZapFields(keysAndValues...)
	z.logger.Error(msg, fields...)
}

// Zap Fatal wrapper
func (z *zapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	fields := convertToZapFields(keysAndValues...)
	z.logger.Fatal(msg, fields...)
}

// Zap With wrapper
func (z *zapLogger) With(keysAndValues ...interface{}) Logger {
	fields := convertToZapFields(keysAndValues...)
	return &zapLogger{logger: z.logger.With(fields...)}
}

// Sync flushes any buffered log entries
func (z *zapLogger) Sync() error {
	return z.logger.Sync()
}

// Helper converters to create zap.Field values that can be passed through
// the public helper methods. We return zap.Field as interface{} so it
// satisfies the Logger interface contract.

func (z *zapLogger) ToString(key, value string) interface{} {
	return zap.String(key, value)
}

func (z *zapLogger) ToDuration(key string, value interface{}) interface{} {
	if d, ok := value.(time.Duration); ok {
		return zap.Duration(key, d)
	}
	// fallback to any
	return zap.Any(key, value)
}

func (z *zapLogger) ToInt(key string, value int) interface{} {
	return zap.Int(key, value)
}

func (z *zapLogger) ToAny(key string, value interface{}) interface{} {
	return zap.Any(key, value)
}

func (z *zapLogger) ToError(err error) interface{} {
	return zap.Error(err)
}

// convertToZapFields converts a mixed slice of key/value pairs and zap.Field
// values into a slice of zap.Field for use with the non-sugared zap.Logger.
func convertToZapFields(kv ...interface{}) []zap.Field {
	var fields []zap.Field
	i := 0
	for i < len(kv) {
		switch v := kv[i].(type) {
		case zap.Field:
			fields = append(fields, v)
			i++
		default:
			// treat as key if next element exists
			if i+1 < len(kv) {
				key := fmt.Sprint(kv[i])
				val := kv[i+1]
				fields = append(fields, zap.Any(key, val))
				i += 2
			} else {
				// lone value, store under generic key
				fields = append(fields, zap.Any(fmt.Sprintf("arg%d", i), kv[i]))
				i++
			}
		}
	}
	return fields
}
