package logd

import (
	"bytes"
	"errors"
	"testing"
)

func TestLevelSwitch(t *testing.T) {
	var buff bytes.Buffer
	var dev = TestModeLog("app.Debug", &buff)

	dev.Log("4021", InfoLevel, "LoadConfig", "Configuratio Loaded")

	dev.SwitchLevel(ErrorLevel)

	// this will be ignored if its level is below the current set log level
	dev.Log("4021", InfoLevel, "LoadConfig", "loading app.config file from disk")

	dev.SwitchLevel(DataTraceLevel)
	dev.Log("4021", InfoLevel, "LoadConfig", "loading app.config file errored out", errors.New("File Not Found!"))
}

// TestBasicLogging tests the output response from using the log api
func TestBasicLogging(t *testing.T) {
	var buff bytes.Buffer
	var dev = TestModeLog("app.Debug", &buff)

	ctx := "3432"
	lvl := InfoLevel
	funcName := "CallRouters#300:36"
	Message := "Initializing Routing Stats"

	// dev.SwitchMode(User)
	dev.Log(ctx, lvl, funcName, Message)
	testRes := basicFormatter(dev, ctx, funcName, Message, nil)

	if buff.String() != testRes {
		t.Fatalf("Invalid response with expected output: Expected %s got %s", testRes, buff.String())
	}

	t.Log("Basic Log format passed")
}

// Switch logLevel to DataTrace and send out some data to include in the trace lines
func TestModeLogging(t *testing.T) {
	var buff bytes.Buffer
	var dev = TestModeLog("app.Debug", &buff)

	ctx := "go.4321"
	funcName := "Agg.WriteResponse#30:3"
	Message := "Sending Response Body"

	// ByteFormatter turns bytes into readable format using hex notation:
	/*<Data
			   0x000: 00 01 03 05 10 ...
	  Data>\n\r
	*/

	var bo = ByteFormatter([]byte(`Thunder routers`))

	//test in dev mode first
	dev.Log(ctx, DataTraceLevel, funcName, Message, bo)
	devtestRes := basicFormatter(dev, ctx, funcName, Message+bo.Format(), nil)

	if buff.String() != devtestRes {
		t.Fatalf("Invalid response with expected output in dev mode: Expected %s got %s", devtestRes, buff.String())
	}

	buff.Reset()
	//switch into user mode
	dev.SwitchMode(UserMode)

	dev.Log(ctx, DataTraceLevel, funcName, Message, bo)
	usertestRes := basicFormatter(dev, ctx, funcName, Message+bo.Format(), nil)

	if buff.String() != usertestRes {
		t.Fatalf("Invalid response with expected output in user mode: Expected %s got %s", usertestRes, buff.String())
	}

}
