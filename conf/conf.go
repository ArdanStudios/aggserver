package conf

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// c is the singular map storage for configurations loaded using Init.
var c = make(map[string]string)

// GetString returns the giving value of the key if found and a boolean to indicate
// if the key exists or not
func GetString(key string) (value string, state bool) {
	value, state = c[key]
	return
}

// GetInt returns the giving integer value of the key if found and
// a boolean to indicate if the key exists or not or was not parsable to an int
func GetInt(key string) (value int, state bool) {
	val, ok := GetString(key)

	// if no key was found, return.
	if !ok {
		return
	}

	// convert the string into a int value.
	intv, err := strconv.Atoi(val)

	if err != nil {
		return
	}

	// value was successfully converted.
	value = intv
	state = true

	return
}

// GetTime returns the giving value of the key as a time.Time object,
// if found and a boolean to indicate if the key exists or not or was not
// a valid parsable time string
func GetTime(key string) (value time.Time, state bool) {
	val, ok := GetString(key)

	// did we find a associate key?
	if !ok {
		return
	}

	// parse the value into a UTC format
	ms, err := time.Parse("2015/01/01 12:00:40.400", val)

	if err != nil {
		return
	}

	value = ms
	state = true
	return
}

// Init is to be called only once to load up the giving namespace if found from the environment/exported
//variables into a global map accessible using the getString,getTime,getInt
//functions. All keys will be made into lowerclass using strings.ToLower.
func Init(namespace string) error {
	// this boolean value used to indicate atleast a once matching case with the
	// namespace.
	var found bool
	var target = c

	// get the lists of available environment variables.
	envs := os.Environ()

	if len(envs) == 0 {
		return fmt.Errorf("No environment variables found")
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
		return fmt.Errorf("No environment variable set for %q", namespace)
	}

	return nil
}
