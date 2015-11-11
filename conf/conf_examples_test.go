package conf_test

import (
	"os"
	"testing"
	"time"

	"github.com/ArdanStudios/aggserver/conf"
)

// ExampleDev shows how to use the conf package
func ExampleDev(t *testing.T) {
	// Setting up some basic environment variables.
	os.Setenv("DOCK_IP", "40.23.233.10")
	os.Setenv("DOCK_PORT", "4044")
	os.Setenv("DOCK_InitStamp", time.Now().String())

	// Init must be called once with the given namespace in to load the
	// associated varaibles, namespace should be in lowercase.
	conf.Init("dock")

	// NOTE: All keys must be in lowercase. The second returned value
	// is a boolean, to indicate wether the key was found and if it was parsable
	// to the desired type.

	// use GetString to retrieve the ip
	_, _ = conf.GetString("ip")

	// use GetInt to retrieve the port if found and parsable to int
	_, _ = conf.GetInt("port")

	// use GetTime to retrieve the InitStamp value as a time object if parsable
	// as such
	_, _ = conf.GetTime("initstamp")

}
