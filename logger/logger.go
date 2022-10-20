package logger

import (
	"io"
	"log"

	"github.com/fatih/color"
)

// Logger is a wrapper around the standard logger.
type Logger interface {
	Info(msgFormat string, params ...interface{})
	Warn(msgFormat string, params ...interface{})
	Error(msgFormat string, params ...interface{})
}

type logger struct {
	logger     *log.Logger
	paramColor func(a ...interface{}) string
	infoColor  func(w io.Writer, format string, a ...interface{})
	warnColor  func(w io.Writer, format string, a ...interface{})
	errColor   func(w io.Writer, format string, a ...interface{})
}

// NewLogger creates a new logger.
func NewLogger() Logger {
	return logger{
		logger:     log.Default(),
		paramColor: color.New(color.FgBlue).SprintFunc(),
		infoColor:  color.New(color.FgGreen).FprintfFunc(),
		warnColor:  color.New(color.FgYellow).FprintfFunc(),
		errColor:   color.New(color.FgRed).FprintfFunc(),
	}
}

// Info logs an info message.
func (l logger) Info(msgFormat string, params ...interface{}) {
	var coloredParams []interface{}
	for _, param := range params {
		if val, ok := param.(string); ok {
			coloredParams = append(coloredParams, color.BlueString(val))
		}
	}

	l.infoColor(l.logger.Writer(), msgFormat, coloredParams...)
}

// Warn logs a warning message.
func (l logger) Warn(msgFormat string, params ...interface{}) {
	l.warnColor(l.logger.Writer(), msgFormat, params...)
}

// Error logs an error message.
func (l logger) Error(msgFormat string, params ...interface{}) {
	l.errColor(l.logger.Writer(), msgFormat, params...)
}
