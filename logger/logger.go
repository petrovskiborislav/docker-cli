package logger

import (
	"io"
	"log"

	"github.com/fatih/color"
)

type Logger struct {
	logger     *log.Logger
	paramColor func(a ...interface{}) string
	infoColor  func(w io.Writer, format string, a ...interface{})
	warnColor  func(w io.Writer, format string, a ...interface{})
	errColor   func(w io.Writer, format string, a ...interface{})
}

// NewLogger creates a new logger.
func NewLogger() Logger {
	return Logger{
		logger:     log.Default(),
		paramColor: color.New(color.FgBlue).SprintFunc(),
		infoColor:  color.New(color.FgGreen).FprintfFunc(),
		warnColor:  color.New(color.FgYellow).FprintfFunc(),
		errColor:   color.New(color.FgRed).FprintfFunc(),
	}
}

// Info logs an info message.
func (l Logger) Info(msgFormat string, params ...interface{}) {
	var coloredParams []interface{}
	for _, param := range params {
		if val, ok := param.(string); ok {
			coloredParams = append(coloredParams, color.BlueString(val))
		}
	}

	l.infoColor(l.logger.Writer(), msgFormat, coloredParams...)
}

// Warn logs a warning message.
func (l Logger) Warn(msgFormat string, params ...interface{}) {
	l.warnColor(l.logger.Writer(), msgFormat, params...)
}

// Error logs an error message.
func (l Logger) Error(msgFormat string, params ...interface{}) {
	l.errColor(l.logger.Writer(), msgFormat, params...)
}
