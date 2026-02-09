package interfaces

type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(err error, msg string, keysAndValues ...interface{})
	Fatal(err error, msg string, keysAndValues ...interface{})
	With(keysAndValues ...interface{}) Logger
}
