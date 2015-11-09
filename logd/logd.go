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
	// User mode only requires a simple readable format
	User Mode = iota + 1
	// Dev mode only requires a an extended information regarding output
	Dev
	// NotSupportedMode output modes that have no supported
	NotSupportedMode
)

// logLevelModeAssoc provides a key:formatstring association for log mode
var logLevelModeAssoc = map[Mode]string{
	1: `Type: %s Level: %s Time: %s Context: %s Func: %s Message: %s`,
	2: `Type: %s Level: %s Time: %s Context: %s Func: %s Line: %s Message: %s`,
}

// basicFormatter formats out the output of the log
func basicFormatter(lg *Loggly, ctx interface{}, funcName, funcMeta, Message string, data ...interface{}) string {
	levelName := logLevelAssoc[lg.Level()]
	modeVal := logLevelModeAssoc[lg.Mode()]
	var ms string

	if atomic.LoadInt32(&lg.testMode) == 0 {
		ms = time.Date(2009, time.November, 10, 15, 0, 0, 0, time.UTC).UTC().Format(layout)
	} else {
		ms = time.Now().UTC().Format(layout)
	}

	if lg.Mode() == User {
		return fmt.Sprintf(modeVal, lg.logType, levelName, ctx, ms, funcName, fmt.Sprintf(Message, data...))
	}

	return fmt.Sprintf(modeVal, lg.logType, levelName, ctx, ms, funcName, funcMeta, fmt.Sprintf(Message, data...))
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
func (l *Loggly) Log(ctx interface{}, level LogLevel, funcName, Message string, data ...interface{}) {
	// var format string
	if level >= l.Level() && level < NotSupportedLevel {
		// l.log.Printf(format,ctx,level,)
	}
}

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

// DataTrace dumps down the log message included with a json formatted data sets
func (l *Loggly) DataTrace(ctx interface{}, funcName string, Message string, data interface{}) {
}

// DataTracef dumps down the log message included with a json formatted data sets
func (l *Loggly) DataTracef(ctx interface{}, funcName string, Message string, data interface{}, vals ...interface{}) {
}
