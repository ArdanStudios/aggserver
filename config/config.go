package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// EnvConfig loads ups the giving namespace if found from the environment/exported
// variables into a supplied map.
// NOTE: All keys will be made into lowerclass using strings.ToLower
func EnvConfig(namespace string, target map[string]interface{}) error {

	// boolean value used to indicate atleast once matching value with the
	// namespace was found.
	var found bool

	// get the lists of available environment variables
	envs := os.Environ()

	if len(envs) == 0 {
		return fmt.Errorf("No environment variables found")
	}

	// create the uppercase version to meet the standard {NAMESPACE_} format
	uspace := fmt.Sprintf("%s_", strings.ToUpper(namespace))

	// loop through and match section using the uppercase namespace using
	// strings.HasPrefix
	for _, val := range envs {
		if !strings.HasPrefix(val, uspace) {
			continue
		}

		found = true
		part := strings.Split(val, "=")
		key := strings.ToLower(strings.TrimPrefix(part[0], uspace))
		value := part[1]

		// assert if value is one of the main go types else set it as just a string

		//if its actually a bool, then save it as such
		if bl, err := strconv.ParseBool(value); err == nil {
			target[key] = bl
			continue
		}

		//if its actually a int(using the int64 range), then save it as such
		if ui, err := strconv.ParseInt(value, 10, 0); err == nil {
			target[key] = ui
			continue
		}

		//if its actually a float, then save it as such
		if fl, err := strconv.ParseFloat(value, 64); err == nil {
			target[key] = fl
			continue
		}

		//if its actually a uint(using the uint64 range), then save it as such
		if ul, err := strconv.ParseUint(value, 10, 64); err == nil {
			target[key] = ul
			continue
		}

		// probably just a string, so we save accordingly
		target[key] = value
	}

	if !found {
		return fmt.Errorf("No environment variable set for %q", namespace)
	}

	return nil
}
