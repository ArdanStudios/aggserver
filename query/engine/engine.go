package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Expression represents a query to be used against the aggregation api
type Expression struct {
	Collection string   `yaml:"collection" json:"collection"`
	Conditions []string `yaml:"conditions" json:"conditions"`
}

// Connect provides a basic connection config struct for creating the necessary
// mongodb connection object.
type Connect struct {
	Host     string
	AuthDb   string
	Username string
	Pass     string
	Db       string
}

// Engine provides a basic mongodb query controller, ensuring concurrent request
// to the given database instance.
type Engine struct {
	*Connect
	*mgo.Session
	wg sync.WaitGroup
}

// New returns a new query engine instance.
func New(c Connect) (*Engine, error) {
	info := &mgo.DialInfo{
		Addrs:    []string{c.Host},
		Timeout:  60 * time.Second,
		Database: c.AuthDb,
		Username: c.Username,
		Password: c.Pass,
	}

	session, err := mgo.DialWithInfo(info)

	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)

	eng := Engine{
		Connect: &c,
		Session: session,
	}

	return &eng, nil
}

// Close ends the session and frees the resources but ensures all
// request are fully processed before ending.
func (e *Engine) Close() {
	e.wg.Wait()
	e.Session.Close()
}

// ResultCallback provides a function type for the return of a result from a
// call to Engine.Query.
type ResultCallback func(error, []bson.M)

func (e *Engine) QueryFile(file string, params map[string]interface{}, rx ResultCallback) {

}

// Query runs a given query []byte slice (contain json strings) against the mongo
// collection using the mgo.Pipe function.
// The query string(json) will be processed into an Expression and a map of
// possible key:value parameters if provided will be used to swap specific pieces
// with the desired values. During the process if any error is encountered, the
// given callback is called and the process returns and frees up the session and
// wait counter
func (e *Engine) Query(query []byte, params map[string]interface{}, rx ResultCallback) {
	go func() {
		//ensure we decrement the wait counter.
		defer e.wg.Done()

		//map out the given expression rules.
		exp, err := mapQuery2Expressions(query, params)

		// if there was an error converting query into a Expression, then reply
		// the callback and return.
		if err != nil {
			rx(err, nil)
			return
		}

		// map out the expression condition lists into a bson.M map.
		evalExpr, err := mapExpressionToBSON(exp)

		// if there was an error converting the lists into bson.M maps, then reply
		// the callback and return.
		if err != nil {
			rx(err, nil)
			return
		}

		scl := e.Copy()
		defer scl.Close()

		// this will contain the result received from the execution of the expression's
		// conditions.
		var result []bson.M

		//get the needed collection from the db.
		col := e.DB(e.Db).C(exp.Collection)

		// evaluate the expression set and get all results associated. If we receive
		// an error, call the callback, reply and return
		if err := col.Pipe(evalExpr).All(&result); err != nil {
			rx(err, nil)
			return
		}

		// no error occured at this point,call the callback with the given result.
		rx(nil, result)
	}()
}

// mapQuery2Expressions will map out a valid json string through a json decoder
// into a Expression struct. If the param map is supplied, it will substitute the
// appropriate key with the value in the map using the following format '#key#'.
func mapQuery2Expressions(query []byte, params map[string]interface{}) (*Expression, error) {
	var exp = &Expression{}

	if err := json.NewDecoder(bytes.NewBuffer(query)).Decode(exp); err != nil {
		return nil, err
	}

	if params != nil && len(params) != 0 {
		for key, value := range params {
			findKey := fmt.Sprintf("#%s#", key)
			useVal := fmt.Sprintf("%q", value)
			for n, rule := range exp.Conditions {
				exp.Conditions[n] = strings.Replace(rule, findKey, useVal, -1)
			}
		}
	}

	return exp, nil
}

// mapExpressionToBSON will take an expressions object and creates an equivalent
// bson.M map for use with a mongodb execution call.
func mapExpressionToBSON(exp *Expression) ([]bson.M, error) {
	var conditions []bson.M

	for _, cond := range exp.Conditions {
		var m = make(bson.M)
		if err := json.NewDecoder(bytes.NewBufferString(cond)).Decode(&m); err != nil {
			return nil, err
		}
		conditions = append(conditions, m)
	}

	return conditions, nil
}
