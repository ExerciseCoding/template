package cuslog

import (
	"io"
	"os"
)

const (
	FmtEmptySeparate = ""
)

// log level
type Level uint8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in production
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

var LevelNameMapping = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel: "INFO",
	WarnLevel: "WARN",
	ErrorLevel: "ERROR",
	PanicLevel: "Panic",
	FatalLevel: "FATAL",
}

// log options
type options struct {
	// 输出位置
	output 			io.Writer
	// 输出位置
	level    		Level
	stdLevel 		Level
	// 输出格式
	formatter 		Formatter
	// 是否开启文件名和行号
	disableCaller 	bool
}

type Option func(*options) 

func initOptions(opts ...Option) (o *options) {
	o = &options{}
	for _, opt := range opts {
		opt(o)
	}

	if o.output == nil {
		o.output = os.Stderr
	}

	if o.formatter == nil {
		o.formatter = &TextFormatter{}
	}

	return
}

func WithOutput(output io.Writer) Option {
	return func(o *options) {
		o.output = output
	}
}

func WithLevel(level Level) Option {
	return func(o *options) {
		o.level = level
	}
}

func WithStdLevel(level Level) Option {
	return func(o *options) {
		o.stdLevel = level
	}
}

func WithFormatter(formatter Formatter) Option {
	return func(o *options) {
		o.formatter = formatter
	}
}

func WithDisableCaller(caller bool) Option {
	return func(o *options) {
		o.disableCaller = caller
	}
}