package log

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

// succeed is the Unicode codepoint for a check mark.
const succeed = "\u2713"

// failed is the Unicode codepoint for an X mark.
const failed = "\u2717"

var doOne = new(sync.Once)
var buf = new(bytes.Buffer)

func initTests() {
	doOne.Do(func() {
		Init(buf, func() int { return DEV })
	})
	SwitchLevel(DEV)
	resetBuffer()
}

func resetBuffer() {
	buf.Reset()
}

func TestLogLevels(t *testing.T) {
	initTests()
	t.Log("Given the log is initialized.")
	{
		t.Log("When we log DEV level messages")
		{
			resetBuffer()

			Dev("5312", "LogWatch", "Error occured retrieving username %s",
				errors.New("Offline Connection!"))

			if !strings.Contains(buf.String(), "Context: 5312") {
				fail(t, "Should contain context: 4312")
			}
			pass(t, "Should contain context: 4312")

			if !strings.Contains(buf.String(), "Func: LogWatch") {
				fail(t, "Should contain 'Func: LogWatch'")
			}
			pass(t, "Should contain 'Func: LogWatch'")

			components := strings.Split(buf.String(), ",")

			if len(components) < 7 {
				fail(t, "Should contain 7 components parts")
			}
			pass(t, "Should contain 7 components parts")
		}

		t.Log("When we must not output DEV messages in USER level")
		{
			resetBuffer()
			SwitchLevel(USER)

			Dev("4312", "RetrieveUser", "Error occured retrieving username %s",
				errors.New("Offline Connection!"))

			if strings.Contains(buf.String(), "Context: 4312") {
				fail(t, "Should have a long that contains context: 4312")
			}
			pass(t, "Should have a long that contains context: 4312")

			if buf.Len() > 0 {
				fail(t, "Should have an empty log when trying to log in User level with Dev()")
			}
			pass(t, "Should have an empty log when trying to log in User level with Dev()")
		}

		t.Log("When we log USER level messages")
		{
			resetBuffer()
			SwitchLevel(USER)

			User("4394", "LogWatch", "Error occured retrieving username %s",
				errors.New("Offline Connection!"))

			if !strings.Contains(buf.String(), "Context: 4394") {
				fail(t, "Should contain context: 4394")
			}
			pass(t, "Should contain context: 4394")

			if strings.Contains(buf.String(), "Func: LogWatch") {
				fail(t, "Should not contain 'Func: LogWatch'")
			}
			pass(t, "Should not contain 'Func: LogWatch'")

			components := strings.Split(buf.String(), ",")

			if len(components) < 5 {
				fail(t, "Should contain 5 components parts")
			}
			pass(t, "Should contain 5 components parts")
		}
	}
}

// TestUserRetrieveLog validates the log trace which is produced when logging
// the response from a UserRetrieve function call.
func TestUserRetrieveLog(t *testing.T) {
	initTests()

	t.Log("Given the need to log the username retrievals.")
	{
		t.Log("When we need to output error status to dev")
		{
			Dev("4312", "RetrieveUser", "Error occured retrieving username %s",
				errors.New("Offline Connection!"))

			if !strings.Contains(buf.String(), "Context: 4312") {
				fail(t, "Should have a log that contains context: 4312")
			}

			pass(t, "Should have a log that contains context: 4312")
		}

		t.Log("When we must output USER messages in DEV level")
		{
			resetBuffer()

			User("4312", "RetrieveUser", "Error occured retrieving username %s",
				errors.New("Offline Connection!"))

			if !strings.Contains(buf.String(), "Context: 4312") {
				fail(t, "Should not have 'context' within log")
			}
			pass(t, "Should not have 'context' within log")

			if buf.Len() == 0 {
				fail(t, "Should not have empty log traces")
			}
			pass(t, "Should not have empty log traces")
		}
	}
}

// TestLogCorrectness validates the logged output content and length when using
// the two available log levels.
func TestLogLines(t *testing.T) {
	initTests()

	t.Log("Given the need to validate log format components.")
	{

		t.Log("When logging in DEV level")
		{

			context := "32"
			funcName := "Munch"
			message := "Raffle Munch SuperBall."
			file := "log_test.go#169"
			curTime := time.Now().UTC().Format(layout)
			pid := os.Getpid()

			expected := fmt.Sprintf(devFormat, pid, "DEV", curTime, context, file, funcName, message+"\n")
			Dev(context, funcName, message)

			if buf.String() != expected {
				fail(t, "Should match log output with expected")
			}
			pass(t, "Should match log output with expected")

		}

		t.Log("When logging in USER level")
		{

			resetBuffer()
			SwitchLevel(USER)

			context := "64"
			funcName := "Munch"
			message := "Raffle Munch SuperBall."
			curTime := time.Now().UTC().Format(layout)
			pid := os.Getpid()

			expected := fmt.Sprintf(userformat, pid, "USER", curTime, context, message+"\n")
			User(context, funcName, message)

			if buf.String() != expected {
				fail(t, "Should match log output with expected")
			}
			pass(t, "Should match log output with expected")

		}
	}

}

// fail is used to log a fail message.
func fail(t *testing.T, message string, data ...interface{}) {
	if len(data) == 0 {
		t.Fatalf("%s. %s", message, failed)
	} else {
		t.Fatalf("%s. %s", fmt.Sprintf(message, data...), failed)
	}
}

// pass is used to log a success message.
func pass(t *testing.T, message string, data ...interface{}) {
	if len(data) == 0 {
		t.Logf("%s. %s", message, succeed)
	} else {
		t.Logf("%s. %s", fmt.Sprintf(message, data...), succeed)
	}
}
