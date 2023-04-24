package logs

var Empty Logger = &emptyLogger{}

type emptyLogger struct{}

func (*emptyLogger) GetLevel() int { return -1 }

func (*emptyLogger) SetLevel(level int) {}

func (*emptyLogger) Fatal(format string, a ...interface{}) {}

func (*emptyLogger) Error(format string, a ...interface{}) {}

func (*emptyLogger) Warn(format string, a ...interface{}) {}

func (*emptyLogger) Info(format string, a ...interface{}) {}

func (*emptyLogger) Debug(format string, a ...interface{}) {}

func (*emptyLogger) Write(level int, format string, a ...interface{}) {}
