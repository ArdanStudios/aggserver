package cfg_test

import (
	"log"
	"os"
	"testing"

	"github.com/ArdanStudios/aggserver/cfg"
)

// succeed is the Unicode codepoint for a check mark.
const succeed = "\u2713"

// t.Fatalfed is the Unicode codepoint for an X mark.
const failed = "\u2717"

// TestLoadingEnvironmentConfig validates the ability to load configuration values
// using the OS-level environment variables.
func TestLoadingEnvironmentConfig(t *testing.T) {
	t.Log("Given a set of environment variables.")
	{
		os.Setenv("MYAPP_PROC_ID", "322")
		os.Setenv("MYAPP_SOCKET", "./tmp/sockets.po")
		os.Setenv("MYAPP_PORT", "4034")

		cfg.Init("myapp")

		t.Log("\tWhen giving a namspace key to search for")
		{

			if cfg.Int("proc_id") != 322 {
				t.Errorf("\t\t%s Should have key %q with value %d", failed, "proc_id", 322)
			} else {
				t.Logf("\t\t%s Should have key %q with value %d", succeed, "proc_id", 322)
			}

			if cfg.String("socket") != "./tmp/sockets.po" {
				t.Errorf("\t\t%s Should have key %q with value %q", failed, "socket", "./tmp/sockets.po")
			} else {
				t.Logf("\t\t%s Should have key %q with value %q", succeed, "socket", "./tmp/sockets.po")
			}

			if cfg.Int("port") != 4034 {
				t.Errorf("\t\t%s Should have key %q with value %d", failed, "port", 4034)
			} else {
				t.Logf("\t\t%s Should have key %q with value %d", succeed, "port", 4034)
			}

		}

		t.Log("\tWhen validating config resposne")
		{

			shouldNotPanic(t, "socket", func() {
				cfg.String("socket")
			})

			shouldNotPanic(t, "proc_id", func() {
				cfg.Int("proc_id")
			})

			shouldPanic(t, "stamp", func() {
				cfg.Time("stamp")
			})

			shouldPanic(t, "pid", func() {
				cfg.Int("pid")
			})

			shouldPanic(t, "dest", func() {
				cfg.String("dest")
			})

		}
	}
}

// shouldPanic receives a context string and a function to run, if the function
// panics, it is considered a success else a failure.
func shouldPanic(t *testing.T, context string, fx func()) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("\t\t%s Should paniced when giving unknown key %q.", failed, context)
		} else {
			t.Logf("\t\t%s Should paniced when giving unknown key %q.", succeed, context)
		}
	}()
	fx()
}

// shouldNotPanic receives a context string and a function to run, if the function
// does not panics, it is considered a success else a failure.
func shouldNotPanic(t *testing.T, context string, fx func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic: %s", err)
			t.Errorf("\t\t%s Should have not paniced when giving unknown key %q.", failed, context)
		} else {
			t.Logf("\t\t%s Should have not paniced when giving unknown key %q.", succeed, context)
		}
	}()
	fx()
}
