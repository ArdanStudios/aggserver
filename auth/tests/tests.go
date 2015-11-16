// Package tests provides tests support and initializations.
package tests

import (
	"bytes"
	"os"
	"testing"

	"github.com/ArdanStudios/aggserver/log"
)

// Succeed is the Unicode codepoint for a check mark.
const Succeed = "\u2713"

// Failed is the Unicode codepoint for an X mark.
const Failed = "\u2717"

// logdest implements io.Writer and is the log package destination.
var logdest bytes.Buffer

// ResetLog can be called at the beginning of a test or example.
func ResetLog() { logdest.Reset() }

// DisplayLog can be called at the end of a test or example.
// It only prints the log contents if the -test.v flag is set.
func DisplayLog() {
	if !testing.Verbose() {
		return
	}
	logdest.WriteTo(os.Stdout)
}

func init() {
	// TODO: Need to read configuration.
	log.Init(&logdest, func() int { return log.DEV })
	// session.Init(nil)
}
