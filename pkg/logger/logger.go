package logger

type Logger interface {
	// Basic logging layers
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})

	// Create derived logger with fields
	With(keysAndValues ...interface{}) Logger

	// Sync flushes any buffered log entries
	Sync() error

	// Helper methods for common field types
	ToString(key, value string) interface{}
	ToDuration(key string, value interface{}) interface{}
	ToInt(key string, value int) interface{}
	ToAny(key string, value interface{}) interface{}
	ToError(err error) interface{}
}