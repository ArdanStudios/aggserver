package main

import (
	"github.com/aggserver/query/engine"
)

func main() {

	engine := engine.New(&engine.Connect{
		Host:     "ds035428.mongolab.com:35428",
		AuthDb:   "goinggo",
		Username: "guest",
		Pass:     "welcome",
		Db:       "goinggo",
	})

	go engine.Serve()
}
