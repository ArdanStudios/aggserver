package testUtils

import (
	"fmt"
	"testing"
)

// succeed is the Unicode codepoint for a check mark.
const succeed = "\u2713"

// failed is the Unicode codepoint for an X mark.
const failed = "\u2717"

// Fail is used to log a fail log message
func Fail(t *testing.T, message string, data ...interface{}) {
	if len(data) == 0 {
		t.Fatalf("%s %s", message, failed)
	} else {
		t.Fatalf("%s %s", fmt.Sprintf(message, data...), failed)
	}
}

// Pass is used to log a success log message
func Pass(t *testing.T, message string, data ...interface{}) {
	if len(data) == 0 {
		t.Logf("%s %s", message, succeed)
	} else {
		t.Logf("%s %s", fmt.Sprintf(message, data...), succeed)
	}
}
