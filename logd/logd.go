package logd

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// UTC Time Layout string
const layout = "2015/04/01 12:00:00.000"

// LogLevel defines a int type represent the different supported loglevels
type LogLevel int

const (
	//NoLogging defines the level where logging is disabled
	NoLogging LogLevel = iota
	// InfoLevel represents the info log level
	InfoLevel
	// DebugLevel represents the debug log level
	DebugLevel
	// DataTraceLevel represents the data dump log level
	DataTraceLevel
	// TraceLevel represents the function trace log level
	TraceLevel
	// ErrorLevel represents the error log level
	ErrorLevel
	// NotSupportedLevel represents level that are not supported
	NotSupportedLevel
)

// association strings with specific log levels
var logLevelAssoc = map[LogLevel]string{
	1: "INFO",
	2: "DEBUG",
	3: "DATATRACE",
	4: "TRACE",
	5: "ERROR",
}

// Mode is used to represent the output format, user log or dev log
type Mode int

const (
	// DevMode mode only requires a an extended information regarding output
	DevMode Mode = iota + 1
	// UserMode mode only requires a simple readable format
	UserMode
	// NotSupportedMode output modes that have no supported
	NotSupportedMode
)

var devFormat = `Type: %s Level: %s Time: %s Context: %s Func: %s Message: %s`
var userFormat = `Type: %s Level: %s Time: %s Context: %s Func: %s Message: %s`

// basicFormatter formats out the output of the log
func basicFormatter(lg *Loggly, ctx interface{}, funcName, Message string, data ...interface{}) string {
	var ms string
	levelName := logLevelAssoc[lg.Level()]

	if atomic.LoadInt32(&lg.testMode) == 1 {
		ms = time.Date(2015, time.November, 10, 15, 0, 0, 0, time.UTC).UTC().Format(layout)
	} else {
		ms = time.Now().UTC().Format(layout)
	}

	if lg.Mode() > DevMode {
		return fmt.Sprintf(userFormat, lg.logType, levelName, ctx, ms, fmt.Sprintf(Message, data...))
	}

	return fmt.Sprintf(devFormat, lg.logType, levelName, ctx, ms, funcName, fmt.Sprintf(Message, data...))
}

// Loggly provides a base logging structure that provides a simple but adequate logging mechanism which provides both human readable and machine readable code
type Loggly struct {
	log      *log.Logger
	logType  string
	ro       sync.RWMutex
	level    LogLevel
	mo       sync.RWMutex
	mode     Mode
	testMode int32
}

// New returns a new instance of Loggly with the currently set loglevel at 1
func New(t string, dev io.Writer) *Loggly {
	lg := Loggly{
		log:     log.New(dev, "", 0),
		logType: t,
		level:   1,
	}
	return &lg
}

// TestModeLog returns a new instance of Loggly with the currently set loglevel at 1
func TestModeLog(t string, dev io.Writer) *Loggly {
	lg := New(t, os.Stdout)
	atomic.StoreInt32(&lg.testMode, 1)
	return lg
}

// StdLog returns a new instance of Loggly with the output device set to stdout
func StdLog(t string) *Loggly {
	return New(t, os.Stdout)
}

// SwitchMode sets the current mode into log instance to the supplied mode instance
func (l *Loggly) SwitchMode(m Mode) {
	//if its not a mode we support, skip
	if m < 0 || m >= NotSupportedMode {
		return
	}
	l.mo.Lock()
	l.mode = m
	l.mo.Unlock()
}

// SwitchLevel sets the current level into log instance
func (l *Loggly) SwitchLevel(lvl LogLevel) {
	//if its not a level we support, skip
	if lvl < 0 || lvl >= NotSupportedLevel {
		return
	}
	l.ro.Lock()
	l.level = lvl
	l.ro.Unlock()
}

// Mode returns the current output mode
func (l *Loggly) Mode() (m Mode) {
	l.mo.RLock()
	m = l.mode
	l.mo.RUnlock()
	return
}

// Level returns the current log level
func (l *Loggly) Level() (lvl LogLevel) {
	l.ro.RLock()
	lvl = l.level
	l.ro.RUnlock()
	return
}

// Log provides the core logging function used by Loggly
func (l *Loggly) Log(ctx interface{}, level LogLevel, funcName string, messages ...interface{}) {
	// var format string
	if level >= l.Level() && level < NotSupportedLevel {
		// l.log.Printf(format,ctx,level,)
	}
}

// Logf provides the core logging function used by Loggly
func (l *Loggly) Logf(ctx interface{}, level LogLevel, funcName, Message string, data ...interface{}) {
}

// Formatter presents an interface where a data value provides its own format directive
type Formatter interface {
	Format() string
}

// ByteFormatter turns a byte into a proper line string of hexdecimal digits
func ByteFormatter(b []byte) Formatter {
	return nil
}

var _log *Loggly

// User logs out into user mode
func User(ctx interface{}, level LogLevel, funcName string, message ...interface{}) {
	_log.SwitchMode(UserMode)
	_log.Log(ctx, level, funcName, message...)
}

// Userf logs out into user mode, will formatting string provided
func Userf(ctx interface{}, level LogLevel, funcName, message string, data ...interface{}) {
	_log.SwitchMode(UserMode)
	_log.Logf(ctx, level, funcName, message, data...)
}

// Dev logs out into dev mode
func Dev(ctx interface{}, level LogLevel, funcName string, message ...interface{}) {
	_log.SwitchMode(DevMode)
	_log.Log(ctx, level, funcName, message...)
}

// Devf logs out into dev mode, will formatting string provided
func Devf(ctx interface{}, level LogLevel, funcName, message string, data ...interface{}) {
	_log.SwitchMode(DevMode)
	_log.Logf(ctx, level, funcName, message, data...)
}
