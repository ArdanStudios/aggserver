package logd

import (
	"errors"
	"testing"
)

var dev = New("app.Debug")

func TestBasicLogging(t *testing.T) {
	dev.Log("3432", InfoLevel, "CallRouters", "Intializing Routing Stats")
	dev.Logf("3432", ErrorLevel, "CallRouters", "Error with data from %s", errors.New("Req: /home Status 404"))
	dev.Logf("3432", DebugLevel, "CallRouters", "Debuging request for Req: %s", "/home")
}

// Switch logLevel to DataTrace and send out some data to include in the trace lines
func TestDataTrace(t *testing.T) {
	dev.DataTrace("go.4321", "Agg.WriteResponse", "Sending Response body", []byte("Thunder routers"))

	//switch out level to a higher priority
	dev.SwitchLevel(ErrorLevel)

	// this log should be ignored as we have entered a high log
	dev.DataTracef("go.4321", "Agg.WriteResponse", "Response Written with Status: %d", nil, 200)
}

func TestErrorLevels(t *testing.T) {
	dev.SwitchLevel(ErrorLevel)
	// all log levels below the current are ignored
	dev.Info("4021", "LoadConfig", "Configuratio Loaded")

	dev.Info("4021", "LoadConfig", "loading app.config file from disk")

	dev.Errorf("4021", "LoadConfig", "loading app.config file errored out", errors.New("File Not Found!"))
}
