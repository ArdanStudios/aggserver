package logd

import (
	"log"
	"os"
)

// LogLevel defines a int type represent the different supported loglevels
type LogLevel string

const (
	// InfoLevel represents the info log level
	InfoLevel LogLevel = "INFO"
	// DebugLevel represents the debug log level
	DebugLevel LogLevel = "DEBUG"
	// DumpLevel represents the data dump log level
	DumpLevel LogLevel = "DUMP"
	// TraceLevel represents the function trace log level
	TraceLevel LogLevel = "TRACE"
	// ErrorLevel represents the error log level
	ErrorLevel LogLevel = "ERROR"
)

// Loggly provides a base logging structure that provides a simple but adequate logging mechanism which provides both human readable and machine readable code
type Loggly struct {
	log *log.Logger
}

// Log provides the core logging function used by Loggly
func (l *Loggly) Log(ctx interface{}, level LogLevel, funcName, Message string, data ...interface{}) {}

// Logf provides the core logging function used by Loggly
func (l *Loggly) Logf(ctx interface{}, level LogLevel, funcName, Message string, data ...interface{}) {
}

// Info logs debug level info
func (l *Loggly) Info(ctx interface{}, funcName, Message string, data ...interface{}) {}

// Infof logs debug level info
func (l *Loggly) Infof(ctx interface{}, funcName, Message string, data ...interface{}) {}

// Error logs debug level info
func (l *Loggly) Error(ctx interface{}, funcName, Message string, data ...interface{}) {}

// Errorf logs debug level info
func (l *Loggly) Errorf(ctx interface{}, funcName, Message string, data ...interface{}) {}

// Debug logs debug level info
func (l *Loggly) Debug(ctx interface{}, funcName, Message string, data ...interface{}) {}

// Debugf logs debug level info
func (l *Loggly) Debugf(ctx interface{}, funcName, Message string, data ...interface{}) {}

// Dump dumps down the log message included with a json formatted data sets
func (l *Loggly) Dump(ctx interface{}, funcName string, jd interface{}, Message string, data ...interface{}) {
}

// Dumpf dumps down the log message included with a json formatted data sets
func (l *Loggly) Dumpf(ctx interface{}, funcName string, jd interface{}, Message string, data ...interface{}) {
}

var central = log.New(os.Stdout, "", 0)

// User provides a loggly logger for handling user reports
var User = Loggly{log: central}

// Dev provides a loggly logger for handling dev reports
var Dev = Loggly{log: central}

// func User(ctx interface{}, level LogLevel, funcName, Message string, data ...interface{}) {
// 	user.Log(ctx, level, funcName, Message, data...)
// }
//
// func Userf(ctx interface{}, level LogLevel, funcName, Message string, data ...interface{}) {
// 	user.Logf(ctx, level, funcName, Message, data...)
// }
//
// func Dev(ctx interface{}, level LogLevel, funcName, Message string, data ...interface{}) {
// 	dev.Log(ctx, level, funcName, Message, data...)
// }
//
// func Devf(ctx interface{}, level LogLevel, funcName, Message string, data ...interface{}) {
// 	dev.Logf(ctx, level, funcName, Message, data...)
// }
