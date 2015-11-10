// Package config provides a package level configuration loader which loaders a
// given set of configuration options using a given namespace and a map as the
// storage endpoint.
// To load a your configuration from your environment using a namespace with the
// following all set within the environment
//		os.Setenv("MYAPP_PROC_ID", "322")
//		os.Setenv("MYAPP_SOCKET", "./tmp/sockets.po")
//		os.Setenv("MYAPP_PORT", "4034")
//
// Simple Do
//      config := make(map[string]interface{})
//      err := EnvConfig("myapp",config)
//
package config
