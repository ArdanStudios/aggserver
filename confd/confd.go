package confd

// Loader provides an interface defining a configuration loader which wraps the process of
// collecting configuration values to be loaded into a provided object pointer,it follows the spirit of
// go serialziation interface style(Marshall/UnMarhsall)
type Loader interface {
	Load(string, interface{}) error
}

// EnvLoader provides a configuration loader which loads straight from the enviroment variables using a supplied namespace
type envLoader struct{}

// EnvConfig provides a global wide environment configuration loader
var EnvConfig = &envLoader{}

// Load expects a string and a map/struct object into which to load the necessarily cofiguration using the namespace giving
func (e *envLoader) Load(namespace string, target interface{}) error {

	return nil
}

// JSONFileLoader provides a configuration loader which loads from a giving file path
type jSONFileLoader struct{}

// JSONConfig provides a generic json configuration loader from a given filepath
var JSONConfig = &jSONFileLoader{}

// Load expects a filepath string and a giving json.UnMarshaller object into which the provided configuration will be loaded
func (e *jSONFileLoader) Load(filePath string, target interface{}) error {

	return nil
}

// YAMLFileLoader provides a configuration loader which loads from a giving file path
type yAMLFileLoader struct{}

// YAMLConfig provides a generic yaml configuration loader from a given filepath
var YAMLConfig = &yAMLFileLoader{}

// Load expects a filepath string and a giving yaml.Unmarhsal interface object into which the provided configuration will be loaded
func (e *yAMLFileLoader) Load(filePath string, target interface{}) error {

	return nil
}
