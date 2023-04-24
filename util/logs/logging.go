package logs

import "github.com/modern-go/reflect2"

type LoggerProvider interface {
	// 获取当前对象上设置的 Logger
	// 如果当前对象没有设置，则返回 Empty
	GetLogger() Logger
	// 为当前对象设置 Logger
	// 如果 logger 为 nil，则会将其置为 Empty
	SetLogger(logger Logger)
}

type LoggingStruct struct {
	logger Logger
}

// 获取当前对象上设置的 Logger
// 如果当前对象没有设置，则返回 Empty
func (this *LoggingStruct) GetLogger() Logger {
	if this != nil && this.logger != nil {
		return this.logger
	}
	return Empty
}

// 为当前对象设置 Logger
// 如果 logger 为 nil，则会将其置为 Empty
func (this *LoggingStruct) SetLogger(logger Logger) {
	if !reflect2.IsNil(logger) {
		this.logger = logger
	} else {
		this.logger = Empty
	}
}

type DefaultLogging struct {
	logger Logger
}

// 获取当前对象上设置的 Logger
// 如果当前对象没有设置，则返回 Empty
func (this *DefaultLogging) GetLogger() Logger {
	if this != nil && this.logger != nil {
		return this.logger
	}
	return Empty
}

// 为当前对象设置 Logger
// 如果 logger 为 nil，则会将其置为 Empty
func (this *DefaultLogging) SetLogger(logger Logger) {
	this.logger = logger
}
