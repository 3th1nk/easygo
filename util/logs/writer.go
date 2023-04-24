package logs

import "io"

func NewWriter(logger Logger, level int) io.Writer {
	return &loggerWriter{logger: logger, level: level}
}

type loggerWriter struct {
	logger Logger
	level  int
}

func (this *loggerWriter) Write(p []byte) (n int, err error) {
	str := string(p)
	Write(this.logger, this.level, str)
	return len(p), nil
}
