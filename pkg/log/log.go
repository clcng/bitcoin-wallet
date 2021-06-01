// Package log provides a global logger for zerolog.
package log

import (
	"context"
	"io"
	"os"

	"github.com/rs/zerolog"
)

// Logger is the global logger.
var Logger zerolog.Logger

func init() {

	//set the zerolog keys map to our filebeat setting
	zerolog.ErrorFieldName = "err"   //don't conflict with filebeat error field
	zerolog.MessageFieldName = "msg" //make it compatible with logrus tag
	zerolog.CallerSkipFrameCount = 3
	// zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// console := zerolog.ConsoleWriter{Out: os.Stdout}
	Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Set as standard logger output
	// stdlog.SetFlags(0)
	// stdlog.SetOutput(Logger)
}

// Output duplicates the global logger and sets w as its output.
func Output(w io.Writer) zerolog.Logger {
	Logger = Logger.Output(w)
	return Logger
}

// Level crestes a child logger with the minium accepted level set to level.
func Level(level zerolog.Level) zerolog.Logger {
	Logger = Logger.Level(level)
	return Logger
}

// With creates a child logger with the field added to its context.
func With() zerolog.Context {
	return Logger.With()
}

// Sample returns a logger with the s sampler.
func Sample(s zerolog.Sampler) zerolog.Logger {
	return Logger.Sample(s)
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func Debug() *zerolog.Event {
	return Logger.Debug()
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func Info() *zerolog.Event {
	return Logger.Info()
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func Warn() *zerolog.Event {
	return Logger.Warn()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func Error() *zerolog.Event {
	return Logger.Error().Caller()
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method.
//
// You must call Msg on the returned event in order to send the event.
func Fatal() *zerolog.Event {
	return Logger.Fatal()
}

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call Msg on the returned event in order to send the event.
func Panic() *zerolog.Event {
	return Logger.Panic()
}

// Log starts a new message with no level. Setting zerolog.GlobalLevel to
// zerlog.Disabled will still disable events produced by this method.
//
// You must call Msg on the returned event in order to send the event.
func Log() *zerolog.Event {
	return Logger.Log()
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	Logger.Print(v...)
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	Logger.Printf(format, v...)
}

// Ctx returns the Logger associated with the ctx. If no logger
// is associated, a disabled logger is returned.
func Ctx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}
