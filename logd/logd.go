package logd

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
