package cfg

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// c is the singular map storage for configurations which are loaded by the call
// to Init().
var c = make(map[string]string)

// String returns the value(in type string) of the giving key, else will panic
// if the key  does not exist.
func String(key string) string {
	value, state := c[key]
	if !state {
		panic(fmt.Sprintf("Unknown Key %s !", key))
	}
	return value
}

// Int returns the value(in type int) of the giving key, else will panic
// if the key does not exist or can not be parsed into an int.
func Int(key string) int {
	// get the giving key and convert the string into a int value.
	intv, err := strconv.Atoi(String(key))

	if err != nil {
		panic(fmt.Sprintf("Key %q values is not a int type", key))
	}

	return intv
}

// Time returns the value(in type time.Time) of the giving key, else will panic
// if the key does not exist, or can not be parsed into a time object.
func Time(key string) time.Time {
	// get the value and parse the value into a UTC format.
	ms, err := time.Parse("2015/01/01 12:00:40.400", String(key))

	if err != nil {
		panic(fmt.Sprintf("Key %q values is not a parsable time string", key))
	}

	return ms
}

// Init is to be called only once, to load up the giving namespace if found,
// in the environment/export variables. All keys will be made into lowerclass
// using strings.ToLower.
func Init(namespace string) {
	// this boolean value used to indicate atleast a once matching case with the
	// namespace.
	var found bool
	var target = c

	// get the lists of available environment variables.
	envs := os.Environ()

	if len(envs) == 0 {
		panic("No environment variables found")
	}

	// create the uppercase version to meet the standard {NAMESPACE_} format.
	uspace := fmt.Sprintf("%s_", strings.ToUpper(namespace))

	// loop through and match each section using the uppercase namespace, by using
	// strings.HasPrefix.
	for _, val := range envs {
		if !strings.HasPrefix(val, uspace) {
			continue
		}

		found = true
		part := strings.Split(val, "=")
		key := strings.ToLower(strings.TrimPrefix(part[0], uspace))
		value := part[1]

		// probably just a string, so we save accordingly
		target[key] = value
	}

	if !found {
		panic(fmt.Sprintf("Namespace %q was not found", namespace))
	}
}
