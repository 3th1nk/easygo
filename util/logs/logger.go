package logs

var Default Logger = Stdout()

type Logger interface {
	GetLevel() int
	SetLevel(level int)

	Fatal(format string, a ...interface{})
	Error(format string, a ...interface{})
	Warn(format string, a ...interface{})
	Info(format string, a ...interface{})
	Debug(format string, a ...interface{})
}

type WritableLogger interface {
	Logger

	Write(level int, format string, a ...interface{})
}

type FieldLogger interface {
	Logger

	With(key string, val interface{}) FieldLogger
	WithMulti(map[string]interface{}) FieldLogger
}
