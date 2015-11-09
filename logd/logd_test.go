package logd

import (
	"bytes"
	"errors"
	"testing"
)

func TestBasicLogging(t *testing.T) {
	var buff bytes.Buffer
	var dev = New("app.Debug", &buff)

	ctx := "3432"
	lvl := InfoLevel
	funcName := "CallRouters"
	funcMeta := "300:36"
	Message := "Initializing Routing Stats"

	// dev.SwitchMode(User)
	dev.Log(ctx, lvl, funcName, Message)
	testRes := basicFormatter(dev, ctx, funcName, funcMeta, Message, nil)

	if buff.String() != testRes {
		t.Fatalf("Invalid response with expected output: Expected %s got %s", testRes, buff.String())
	}

	t.Log("Basic Log format passed")
}

// Switch logLevel to DataTrace and send out some data to include in the trace lines
func TestDataTrace(t *testing.T) {
	var buff bytes.Buffer
	var dev = New("app.Debug", &buff)

	ctx := "go.4321"
	funcName := "Agg.WriteResponse"
	funcMeta := "30:3"
	Message := "Sending Response Body"

	/*
		   Turns bytes into:
			   0x000: 00 01 03 05 10 ...
	*/

	var bo = ByteFormatter([]byte(`Thunder routers`))
	dev.DataTrace(ctx, funcName, Message, bo)
	testRes := basicFormatter(dev, ctx, funcName, funcMeta, Message+bo.Format(), nil)

	if buff.String() != testRes {
		t.Fatalf("Invalid response with expected output: Expected %s got %s", testRes, buff.String())
	}

	//switch out level to a higher priority
	dev.SwitchLevel(ErrorLevel)

	// this log should be ignored as we have entered a high log
	dev.DataTracef("go.4321", "Agg.WriteResponse", "Response Written with Status: %d", nil, 200)
}

func TestErrorLevels(t *testing.T) {
	var buff bytes.Buffer
	var dev = New("app.Debug", &buff)
	dev.SwitchLevel(ErrorLevel)
	// all log levels below the current are ignored
	dev.Info("4021", "LoadConfig", "Configuratio Loaded")

	dev.Info("4021", "LoadConfig", "loading app.config file from disk")

	dev.Error("4021", "LoadConfig", "loading app.config file errored out", errors.New("File Not Found!"))
}
