package main

import (
	"log"

	"github.com/ArdanStudios/aggserver/query/engine"

	"gopkg.in/mgo.v2/bson"
)

func main() {

	servo, err := engine.New(engine.Connect{
		Host:     "ds035428.mongolab.com:35428",
		AuthDb:   "goinggo",
		Username: "guest",
		Pass:     "welcome",
		Db:       "goinggo",
	})

	if err != nil {
		panic(err)
	}

	log.Printf("Transactions initiating")
	err = servo.QueryFile("./queries/transactions.json", map[string]interface{}{
		"userId": "396bc782-6ac6-4183-a671-6e75ca5989a5",
	}, func(err error, expr *engine.Expression, result []bson.M) {
		log.Printf("Error %s and Result %s", err, result)
	})

	if err != nil {
		panic(err)
	}

	log.Printf("Requesting close")
	servo.Close()
}
