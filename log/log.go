package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

// UTC Time Layout string to be used in formating time values.
const layout = "2015/04/01 12:00:00.000"

// Level constants that define the supported usable LogLevel.
const (
	USER int = iota
	DEV
)

// loggly is a package level structure that houses the internal logger be used
// for logging events and notifications with in the available log levels.
type loggly struct {
	logMutex   sync.RWMutex
	log        *log.Logger
	level      int
	levelMutex sync.RWMutex
}

// Level returns the current logLevel of the logger.
func (l *loggly) Level() (lvl int) {
	l.levelMutex.RLock()
	lvl = l.level
	l.levelMutex.RUnlock()
	return
}

// SwitchLevel changes the current logLevel of the logger.
func (l *loggly) SwitchLevel(u int) {
	l.levelMutex.Lock()
	l.level = u
	l.levelMutex.Unlock()
}

// _log maintains a pointer to the single log instance used package wide.
var _log = loggly{
// log: log.New(os.Stdout, "", 0),
}

// Init sets up the necessary logging instance with the appropriate loglevel to
// be used by the logger.
// WARNING: This function should only be called once.
func Init(w io.Writer, fx func() int) {
	// lock the global logger access mutex before writing to the global instance.
	_log.logMutex.Lock()
	{
		_log.log = log.New(w, "", 0)
		if fx != nil {
			_log.SwitchLevel(fx())
		}
	}
	_log.logMutex.Unlock()
}

// SwitchLevel provides a means of switching the current LogLevel used by
// the logger.
func SwitchLevel(lvl int) {
	_log.logMutex.RLock()
	defer _log.logMutex.RUnlock()
	_log.SwitchLevel(lvl)
}

// devFormat is the current format used in creating the DEV level logging output
var devFormat = `Pid: %d, Level: %s, Time: %s, Context: %s, File: %s, Func: %s, Message: %s`

// Dev logs trace information out when using the DEV log level.
func Dev(context interface{}, funcName string, format string, a ...interface{}) {

	_log.logMutex.RLock()
	defer _log.logMutex.RUnlock()

	// check the current log level being used,
	// then ensure we are in Dev level to allow logging for DEV.
	// if the loglevel defers, simple return.
	if _log.Level() >= DEV {
		fn, file, pid := fileQuery(funcName, 2)
		_log.log.Printf(devFormat, pid, "DEV", time.Now().UTC().Format(layout), context, file, fn, fmt.Sprintf(format, a...))
	}
}

// userFormat is the current format used in creating the USER level logging output
var userformat = `Pid: %d, Level: %s, Time: %s, Context: %s, Message: %s`

// User logs trace information out when using the USER log level.
func User(context interface{}, funcName string, format string, a ...interface{}) {
	_log.logMutex.RLock()
	defer _log.logMutex.RUnlock()

	// check the current log level being used,
	// then ensure we are in Dev level to allow logging for DEV.
	// if the loglevel defers, simple return.
	if _log.Level() >= USER {
		_, _, pid := fileQuery(funcName, 2)
		_log.log.Printf(userformat, pid, "USER", time.Now().UTC().Format(layout), context, fmt.Sprintf(format, a...))
	}
}

// fileQuery retrieves the functionName(if provided is empty) and line number,
// filePath, and current pid of the running process
func fileQuery(funcName string, depth int) (function string, file string, pid int) {
	pid = os.Getpid()

	if funcName == "" {
		pc := make([]uintptr, depth+1)
		runtime.Callers(depth, pc)
		f := runtime.FuncForPC(pc[depth-1])
		_, function = path.Split(f.Name())
	} else {
		function = funcName
	}

	_, filePath, line, ok := runtime.Caller(depth)

	if !ok {
		file = "unknown.go#0"
		return
	}

	_, file = path.Split(filePath)

	file = fmt.Sprintf("%s#%d", file, line)
	return
}
