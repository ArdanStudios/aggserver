package confd

import (
	"errors"
	"os"
	"testing"

	"github.com/ArdanStudios/aggserver/testUtils"
)

var ErrbadLoadType = errors.New("Invalid type to load into, expected map pointer")

type mapTestLoader string

func (t mapTestLoader) Load(name string, m interface{}) error {
	mx, ok := m.(map[string]interface{})

	if !ok {
		return ErrbadLoadType
	}

	mx["sound"] = name

	return nil
}

// TestConfigLoader tests the use of the Loader interface to create custom configuration loaders
func TestConfigLoader(t *testing.T) {
	var loader mapTestLoader
	var res = map[string]interface{}{}

	if err := loader.Load("duck! duck!", res); err != nil {
		testUtils.Fail(t, "Error occured during load process: %s", err)
	}

	if _, ok := res["sound"]; !ok {
		testUtils.Fail(t, "Expected %q key in config", "sound")
	}

	if res["sound"] != "duck! duck!" {
		testUtils.Fail(t, "Expected %q config value but got %q", "duck! duck!", res["sound"])
	}

	testUtils.Pass(t, "Loader successfully loaded new config value")
}

// TestEnvLoader tests loading configuration from environment variables
func TestEnvLoader(t *testing.T) {
	var res = map[string]interface{}{}

	os.Setenv("MYAPP_PROC_ID", "322")
	os.Setenv("MYAPP_SOCKET", "./tmp/sockets.po")
	os.Setenv("MYAPP_PORT", "4034")

	if err := EnvConfig.Load("myapp", res); err != nil {
		testUtils.Fail(t, "Unable to load environmental config for %s: %s", "myapp", err)
	}

	if _, ok := res["proc_id"]; !ok {
		testUtils.Fail(t, "Expected to find Key %q in map", "proc_id")
	}

	if _, ok := res["socket"]; !ok {
		testUtils.Fail(t, "Expected to find Key %q in map", "socket")
	}

	if _, ok := res["port"]; !ok {
		testUtils.Fail(t, "Expected to find Key %q in map", "port")
	}

	testUtils.Pass(t, "Successfully loaded config from environment variables")
}

type marshalMap map[string]interface{}

func (j *marshalMap) UnmarshalJSON(bo []byte) error {
	//do your own thing here
  return nil
}

func (j *marshalMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	//do your own thing here
  return nil
}

// TestJSONLoader tests loading configuration from a given json file
func TestJSONLoader(t *testing.T) {
	var res = make(marshalMap)

	if err := JSONConfig.Load("myapp.json", &res); err != nil {
		testUtils.Fail(t, "Unable to load json config for %s: %s", "myapp.json", err)
	}

	if _, ok := res["proc_id"]; !ok {
		testUtils.Fail(t, "Expected to find Key %q in map", "proc_id")
	}

	if _, ok := res["socket"]; !ok {
		testUtils.Fail(t, "Expected to find Key %q in map", "socket")
	}

	if _, ok := res["port"]; !ok {
		testUtils.Fail(t, "Expected to find Key %q in map", "port")
	}

	testUtils.Pass(t, "Successfully loaded config from environment variables")
}

// TestYAMLLoader tests loading configuration from a given yaml file
func TestYAMLLoader(t *testing.T) {
	var res = make(marshalMap)

	if err := YAMLConfig.Load("myapp.yml", &res); err != nil {
		testUtils.Fail(t, "Unable to load json config for %s: %s", "myapp.json", err)
	}

	if _, ok := res["proc_id"]; !ok {
		testUtils.Fail(t, "Expected to find Key %q in map", "proc_id")
	}

	if _, ok := res["socket"]; !ok {
		testUtils.Fail(t, "Expected to find Key %q in map", "socket")
	}

	if _, ok := res["port"]; !ok {
		testUtils.Fail(t, "Expected to find Key %q in map", "port")
	}

	testUtils.Pass(t, "Successfully loaded config from environment variables")
}
