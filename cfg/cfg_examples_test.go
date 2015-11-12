package cfg_test

import (
	"fmt"
	"os"
	"time"

	"github.com/ArdanStudios/aggserver/cfg"
)

// ExampleDev shows how to use the conf package
func ExampleDev() {
	// Setting up some basic environment variables.
	os.Setenv("DOCK_IP", "40.23.233.10")
	os.Setenv("DOCK_PORT", "4044")

	ms := time.Date(2009, time.November, 10, 15, 0, 0, 0, time.UTC).UTC().Format(time.UnixDate)
	os.Setenv("DOCK_InitStamp", ms)

	// Init() must be called only once, with the given namespace to load.
	cfg.Init("dock")

	// NOTE: All keys must be in lowercase. The second returned value
	// is a boolean, to indicate wether the key was found and if it was parsable
	// to the desired type.

	// To get the ip string, use cfg.String(key)
	fmt.Println(cfg.String("ip"))

	// To get the port number, use cfg.Int(key)
	fmt.Println(cfg.Int("port"))

	// To get the timestap, use cfg.Time(key)
	fmt.Println(cfg.Time("initstamp"))

	// Output:
	// 40.23.233.10
	// 4044
	// 2009-11-10 15:00:00 +0000 UTC
}
