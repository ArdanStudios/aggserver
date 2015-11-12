package cfg

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// c represents the global struct, with a map used to store the loaded keys
// from the environment variable.
var c struct {
	m map[string]string
}

// Init is to be called only once, to load up the giving namespace if found,
// in the environment/export variables. All keys will be made into lowerclass
// using strings.ToLower.
func Init(namespace string) {
	// Initialize the internal struct map
	c.m = make(map[string]string)

	// This boolean value used to indicate atleast a once matching case with the
	// namespace.
	var found bool

	// Get the lists of available environment variables.
	envs := os.Environ()

	if len(envs) == 0 {
		panic("No environment variables found")
	}

	// Create the uppercase version to meet the standard {NAMESPACE_} format.
	uspace := fmt.Sprintf("%s_", strings.ToUpper(namespace))

	// Loop through and match each section using the uppercase namespace, by using
	// strings.HasPrefix.
	for _, val := range envs {
		if !strings.HasPrefix(val, uspace) {
			continue
		}

		found = true
		part := strings.Split(val, "=")
		key := strings.ToLower(strings.TrimPrefix(part[0], uspace))
		value := part[1]

		// Probably just a string, so we save accordingly.
		c.m[key] = value
	}

	if !found {
		panic(fmt.Sprintf("Namespace %q was not found", namespace))
	}
}

// String returns the value(in type string) of the giving key, else will panic
// if the key does not exist.
func String(key string) string {
	// Get the keys value and bool check if it exists
	value, state := c.m[key]
	// If key does not exists, panic.
	if !state {
		panic(fmt.Sprintf("Unknown Key %s !", key))
	}

	return value
}

// Int returns the value(in type int) of the giving key, else will panic
// if the key does not exist or can not be parsed into an int.
func Int(key string) int {
	// Get the giving key and convert the string into a int value.
	intv, err := strconv.Atoi(String(key))
	// If value can not be converted into an int type, then panic.
	if err != nil {
		panic(fmt.Sprintf("Key %q values is not a int type", key))
	}

	return intv
}

// Time returns the value(in type time.Time) of the giving key, else will panic
// if the key does not exist, or can not be parsed into a time object.
func Time(key string) time.Time {
	// Get the value and attempt to parse it into a time object.
	ms, err := time.Parse(time.UnixDate, String(key))
	// If error occured trying to pass value, then panic.
	if err != nil {
		panic(fmt.Sprintf("%q is unparsable as a time string due to error %s", String(key), err.Error()))
	}

	return ms
}
