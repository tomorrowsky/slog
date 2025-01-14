/*
Package slog Lightweight, extensible, configurable logging library written in Go.

Source code and other details for the project are available at GitHub:

	https://github.com/gookit/slog

Quick usage:

	package main

	import (
		"github.com/gookit/slog"
	)

	func main() {
		slog.Info("info log message")
		slog.Warn("warning log message")
		slog.Infof("info log %s", "message")
		slog.Debugf("debug %s", "message")
	}

More usage please see README.

*/
package slog

import (
	"io"
	"os"
	"time"

	"github.com/gookit/color"
)

// SugaredLogger definition.
// Is a fast and usable Logger, which already contains the default formatting and handling capabilities
type SugaredLogger struct {
	*Logger
	// Formatter log message formatter. default use TextFormatter
	Formatter Formatter
	// Output output writer
	Output io.Writer
	// Level for log handling.
	// Greater than or equal to this level will be recorded
	Level Level
}

// NewSugaredLogger create new SugaredLogger
func NewSugaredLogger(output io.Writer, level Level) *SugaredLogger {
	sl := &SugaredLogger{
		Level:  level,
		Output: output,
		Logger: New(),
		// default value
		Formatter: NewTextFormatter(),
	}

	// NOTICE: use self as an log handler
	sl.AddHandler(sl)

	return sl
}

// NewJSONSugared create new SugaredLogger with JSONFormatter
func NewJSONSugared(out io.Writer, level Level) *SugaredLogger {
	sl := NewSugaredLogger(out, level)
	sl.Formatter = NewJSONFormatter()

	return sl
}

// Configure current logger
func (sl *SugaredLogger) Configure(fn func(sl *SugaredLogger)) *SugaredLogger {
	fn(sl)
	return sl
}

// Reset the logger
func (sl *SugaredLogger) Reset() {
	*sl = *NewSugaredLogger(os.Stdout, ErrorLevel)

	// reset handlers and processors
	// sl.Logger.Reset()

	// NOTICE: use self as an log handler
	// sl.AddHandler(sl)
	// sl.Formatter = NewTextFormatter()
}

// IsHandling Check if the current level can be handling
func (sl *SugaredLogger) IsHandling(level Level) bool {
	return sl.Level.ShouldHandling(level)
}

// Handle log record
func (sl *SugaredLogger) Handle(record *Record) error {
	bts, err := sl.Formatter.Format(record)
	if err != nil {
		return err
	}

	_, err = sl.Output.Write(bts)
	return err
}

// Close all log handlers
func (sl *SugaredLogger) Close() error {
	sl.Logger.VisitAll(func(handler Handler) error {
		if _, ok := handler.(*SugaredLogger); !ok {
			_ = handler.Close()
		}
		return nil
	})
	return nil
}

// FlushAll all logs. alias of the FlushAll()
func (sl *SugaredLogger) Flush() error {
	return sl.FlushAll()
}

// Flush all logs
func (sl *SugaredLogger) FlushAll() error {
	sl.Logger.VisitAll(func(handler Handler) error {
		if _, ok := handler.(*SugaredLogger); !ok {
			_ = handler.Flush()
		}
		return nil
	})
	return nil
}

//
// ------------------------------------------------------------
// Global std logger operate
// ------------------------------------------------------------
//

// std logger is an SugaredLogger.
// It is directly available without any additional configuration
var std = NewStdLogger()

// NewStdLogger instance
func NewStdLogger() *SugaredLogger {
	return NewSugaredLogger(os.Stdout, DebugLevel).Configure(func(sl *SugaredLogger) {
		sl.SetName("stdLogger")
		sl.ReportCaller = true
		// auto enable console color
		sl.Formatter.(*TextFormatter).EnableColor = color.SupportColor()
	})
}

// Std get std logger
func Std() *SugaredLogger {
	return std
}

// Reset the std logger
func Reset() {
	ResetExitHandlers(true)
	// new std
	std = NewStdLogger()
}

// Configure the std logger
func Configure(fn func(logger *SugaredLogger)) {
	std.Configure(fn)
}

// Exit runs all the logger exit handlers and then terminates the program using os.Exit(code)
func Exit(code int) {
	std.Exit(code)
}

// SetExitFunc to the std logger
func SetExitFunc(fn func(code int)) {
	std.ExitFunc = fn
}

// Flush log messages
func Flush() error {
	return std.Flush()
}

// FlushTimeout flush logs with timeout.
func FlushTimeout(timeout time.Duration) {
	std.FlushTimeout(timeout)
}

// FlushDaemon run flush handle on daemon
//
// Usage:
// 	go slog.FlushDaemon()
func FlushDaemon() {
	std.FlushDaemon()
}

// SetLogLevel for the std logger
func SetLogLevel(l Level) {
	std.Level = l
}

// SetFormatter to std logger
func SetFormatter(f Formatter) {
	std.Formatter = f
}

// GetFormatter of the std logger
func GetFormatter() Formatter {
	return std.Formatter
}

// AddHandler to the std logger
func AddHandler(h Handler) {
	std.AddHandler(h)
}

// PushHandler to the std logger
func PushHandler(h Handler) {
	std.AddHandler(h)
}

// AddHandlers to the std logger
func AddHandlers(hs ...Handler) {
	std.AddHandlers(hs...)
}

// PushHandlers to the std logger
func PushHandlers(hs ...Handler) {
	std.PushHandlers(hs...)
}

// AddProcessor to the logger
func AddProcessor(p Processor) {
	std.AddProcessor(p)
}

// AddProcessors to the logger
func AddProcessors(ps ...Processor) {
	std.AddProcessors(ps...)
}

// -------------------------- New record with log data, fields -----------------------------

// WithData new record with data
func WithData(data M) *Record {
	return std.WithData(data)
}

// WithFields new record with fields
func WithFields(fields M) *Record {
	return std.WithFields(fields)
}

// -------------------------- Add log messages with level -----------------------------

// Print logs a message at level PrintLevel
func Print(args ...interface{}) {
	std.Log(PrintLevel, args...)
}

// Println logs a message at level PrintLevel
func Println(args ...interface{}) {
	std.Log(PrintLevel, args...)
}

// Printf logs a message at level PrintLevel
func Printf(format string, args ...interface{}) {
	std.Logf(PrintLevel, format, args...)
}

// Trace logs a message at level Trace
func Trace(args ...interface{}) {
	std.Log(TraceLevel, args...)
}

// Tracef logs a message at level Trace
func Tracef(format string, args ...interface{}) {
	std.Logf(TraceLevel, format, args...)
}

// Info logs a message at level Info
func Info(args ...interface{}) {
	std.Log(InfoLevel, args...)
}

// Infof logs a message at level Info
func Infof(format string, args ...interface{}) {
	std.Logf(InfoLevel, format, args...)
}

// Notice logs a message at level Notice
func Notice(args ...interface{}) {
	std.Log(NoticeLevel, args...)
}

// Noticef logs a message at level Notice
func Noticef(format string, args ...interface{}) {
	std.Logf(NoticeLevel, format, args...)
}

// Warn logs a message at level Warn
func Warn(args ...interface{}) {
	std.Log(WarnLevel, args...)
}

// Warnf logs a message at level Warn
func Warnf(format string, args ...interface{}) {
	std.Logf(WarnLevel, format, args...)
}

// Error logs a message at level Error
func Error(args ...interface{}) {
	std.Log(ErrorLevel, args...)
}

// ErrorT logs a error type at level Error
func ErrorT(err error) {
	if err != nil {
		std.Log(ErrorLevel, err)
	}
}

// Errorf logs a message at level Error
func Errorf(format string, args ...interface{}) {
	std.Logf(ErrorLevel, format, args...)
}

// Debug logs a message at level Debug
func Debug(args ...interface{}) {
	std.Log(DebugLevel, args...)
}

// Debugf logs a message at level Debug
func Debugf(format string, args ...interface{}) {
	std.Logf(DebugLevel, format, args...)
}

// Fatal logs a message at level Fatal
func Fatal(args ...interface{}) {
	std.Log(FatalLevel, args...)
}

// Fatalf logs a message at level Fatal
func Fatalf(format string, args ...interface{}) {
	std.Logf(FatalLevel, format, args...)
}

// Panic logs a message at level Panic
func Panic(args ...interface{}) {
	std.Log(PanicLevel, args...)
}

// Panicf logs a message at level Panic
func Panicf(format string, args ...interface{}) {
	std.Logf(PanicLevel, format, args...)
}
