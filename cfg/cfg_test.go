package cfg_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/ArdanStudios/aggserver/cfg"
)

// succeed is the Unicode codepoint for a check mark.
const succeed = "\u2713"

// failed is the Unicode codepoint for an X mark.
const failed = "\u2717"

// TestLoadingEnvironmentConfig validates the ability to load configuration values
// using the OS-level environment variables.
func TestLoadingEnvironmentConfig(t *testing.T) {
	t.Log("Given a set of environment variables.")
	{
		os.Setenv("MYAPP_PROC_ID", "322")
		os.Setenv("MYAPP_SOCKET", "./tmp/sockets.po")
		os.Setenv("MYAPP_PORT", "4034")

		t.Log(" MYAPP_PROC_ID=322")
		t.Log(" MYAPP_SOCKET='./tmp/sockets.po'")
		t.Log(" MYAPP_PORT=4034")

		t.Log("\tWhen giving a namspace key to search for")
		{

			cfg.Init("myapp")
			pass(t, "\t\tShould load environment variables for namespace %q", "myapp")

			if cfg.Int("proc_id") != 322 {
				fail(t, "\t\tShould have key %q with int type value %d", "proc_id", 322)
			}
			pass(t, "\t\tShould have key %q with int type value %d", "proc_id", 322)

			if cfg.String("socket") != "./tmp/sockets.po" {
				fail(t, "\t\tShould have key %q with string type value %q", "socket", "./tmp/sockets.po")
			}
			pass(t, "\t\tShould have key %q with string type value %q", "socket", "./tmp/sockets.po")

			if cfg.Int("port") != 4034 {
				fail(t, "\t\tShould have key %q with int type value %d", "port", 4034)
			}
			pass(t, "\t\tShould have key %q with int type value %d", "port", 4034)

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
