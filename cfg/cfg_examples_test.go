package cfg_test

import (
	"os"
	"testing"
	"time"

	"github.com/ArdanStudios/aggserver/cfg"
)

// ExampleDev shows how to use the conf package
func ExampleDev(t *testing.T) {
	// Setting up some basic environment variables.
	os.Setenv("DOCK_IP", "40.23.233.10")
	os.Setenv("DOCK_PORT", "4044")
	os.Setenv("DOCK_InitStamp", time.Now().String())

	// Init() must be called only once, with the given namespace to load.
	cfg.Init("dock")

	// NOTE: All keys must be in lowercase. The second returned value
	// is a boolean, to indicate wether the key was found and if it was parsable
	// to the desired type.

	// use String() to retrieve the ip
	_ = cfg.String("ip")

	// use Int() to retrieve the port if found and parsable to int
	_ = cfg.Int("port")

	// use Time() to retrieve the InitStamp value as a time object if parsable
	// as such
	_ = cfg.Time("initstamp")

}
