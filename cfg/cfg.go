package cfg

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// c represents the configuration store, with a map to store the loaded keys
// from the environment.
var c struct {
	m map[string]string
}

// Init is to be called only once, to load up the giving namespace if found,
// in the environment variables. All keys will be made lowerclassed.
func Init(namespace string) {
	// Initialize the internal struct map.
	c.m = make(map[string]string)

	// This boolean value used to indicate atleast that one key was found.
	var foundOne bool

	// Get the lists of available environment variables.
	envs := os.Environ()
	if len(envs) == 0 {
		panic("No environment variables found")
	}

	// Create the uppercase version to meet the standard {NAMESPACE_} format.
	uspace := fmt.Sprintf("%s_", strings.ToUpper(namespace))

	// Loop and match each variable using the uppercase namespace.
	for _, val := range envs {
		if !strings.HasPrefix(val, uspace) {
			continue
		}

		foundOne = true
		part := strings.Split(val, "=")
		c.m[strings.ToLower(strings.TrimPrefix(part[0], uspace))] = part[1]
	}

	if !foundOne {
		panic(fmt.Sprintf("Namespace %q was not found", namespace))
	}
}

// String returns the value(in type string) of the giving key, else will panic
// if the key was not found.
func String(key string) string {
	// Get the key's value and existence status.
	value, state := c.m[key]
	// If key does not exists?, panic.
	if !state {
		panic(fmt.Sprintf("Unknown Key %s !", key))
	}

	return value
}

// Int returns the value(in type int) of the giving key, else panics
// when not found.
func Int(key string) int {
	// Get the key's value and convert to a int.
	intv, err := strconv.Atoi(String(key))
	// If value can not be converted, then panic.
	if err != nil {
		panic(fmt.Sprintf("Key %q values is not a int type", key))
	}

	return intv
}

// Time returns the value(in type time.Time) of the giving key, else panics,
// if the key does not exist, or can not be parsed into a valid time.Time.
func Time(key string) time.Time {
	// Get the value and attempt to parse to time.Time.
	ms, err := time.Parse(time.UnixDate, String(key))
	// If error occured, panic.
	if err != nil {
		panic(fmt.Sprintf("%q is unparsable as a time string due to error %s", String(key), err.Error()))
	}

	return ms
}
