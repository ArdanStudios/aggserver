package config

import (
	"fmt"
	"log"
	"os"
	"testing"
)

// succeed is the Unicode codepoint for a check mark.
const succeed = "\u2713"

// failed is the Unicode codepoint for an X mark.
const failed = "\u2717"

// TestEnvLoader tests loading configuration from environment variables
func TestEnvLoader(t *testing.T) {
	var res = map[string]interface{}{}

	t.Log("Given a set of environment variables")
	{

		os.Setenv("MYAPP_PROC_ID", "322")
		os.Setenv("MYAPP_SOCKET", "./tmp/sockets.po")
		os.Setenv("MYAPP_PORT", "4034")

		t.Log("When giving a namspace key to search for")
		{

			if err := EnvConfig("myapp", res); err != nil {
				fail(t, "Should load environment variables for namespace %q", "myapp")
			}
			pass(t, "Should load environment variables for namespace %q", "myapp")

			if _, ok := res["proc_id"]; !ok {
				fail(t, "Should find key %q in map", "proc_id")
			}
			pass(t, "Should find key %q in map", "proc_id")

			log.Printf("%s", res["proc_id"])

			if res["proc_id"] != int64(322) {
				fail(t, "Should have key %q with value %d in map", "proc_id", 322)
			}
			pass(t, "Should have key %q with value %d in map", "proc_id", 322)

			if _, ok := res["socket"]; !ok {
				fail(t, "Should find key %q in map", "socket")
			}
			pass(t, "Should find key %q in map", "socket")

			if res["socket"] != "./tmp/sockets.po" {
				fail(t, "Should have key %q with value %q in map", "socket", "./tmp/sockets.po")
			}
			pass(t, "Should have key %q with value %q in map", "socket", "./tmp/sockets.po")

			if _, ok := res["port"]; !ok {
				fail(t, "Should find key %q in map", "port")
			}
			pass(t, "Should find key 'port' in map")

			if res["port"] != int64(4034) {
				fail(t, "Should have key %q with value %d in map", "port", 4034)
			}
			pass(t, "Should have key %q with value %d in map", "port", 4034)

		}

	}

}

// fail is used to log a fail message.
func fail(t *testing.T, message string, data ...interface{}) {
	if len(data) == 0 {
		t.Fatalf("%s %s", message, failed)
	} else {
		t.Fatalf("%s %s", fmt.Sprintf(message, data...), failed)
	}
}

// pass is used to log a success message.
func pass(t *testing.T, message string, data ...interface{}) {
	if len(data) == 0 {
		t.Logf("%s %s", message, succeed)
	} else {
		t.Logf("%s %s", fmt.Sprintf(message, data...), succeed)
	}
}
