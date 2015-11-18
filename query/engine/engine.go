package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// session is the global mongodb session for handling requests
var session *mgo.Session
var config *Connect
var wg sync.WaitGroup

// Expression represents a query to be used against the aggregation api
type Expression struct {
	Collection string   `yaml:"collection" json:"collection"`
	Conditions []string `yaml:"conditions" json:"conditions"`
}

// Rule repesent a set of expressions which will be tested against
type Rule struct {
	Test *Expression `json:"test" yaml:"test"`
	Fail *Expression `json:"failed" yaml:"failed"`
	Pass *Expression `json:"passed" yaml:"passed"`
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

// Init must be called only once, it initializes and connects up the session
// and query processor.
func Init(c Connect) {
	info := &mgo.DialInfo{
		Addrs:    []string{c.Host},
		Timeout:  60 * time.Second,
		Database: c.AuthDb,
		Username: c.Username,
		Password: c.Pass,
	}

	ses, err := mgo.DialWithInfo(info)

	if err != nil {
		log.Fatal(err)
	}

	// Reads may not be entirely up-to-date, but they will always see the
	// history of changes moving forward, the data read will be consistent
	// across sequential queries in the same session, and modifications made
	// within the session will be observed in following queries (read-your-writes).
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode
	ses.SetMode(mgo.Monotonic, true)

	session = ses
}

// Wait forces the current routine to wait until all requests are fully processed.
func Wait() {
	wg.Wait()
}

// Session returns a copy of the initialized mognodb session.
func Session() *mgo.Session {
	return session.Copy()
}

// Stringify returns a string representation of a given value,if it could not be
// turned into a string, it will return an empty string.
func Stringify(m interface{}) string {
	val, err := json.Marshal(m)
	if err != nil {
		return ""
	}

	return string(val)
}

// ResultCallback provides a function type for the return of a result from a
// call to Engine.Query.
type ResultCallback func(error, []bson.M)

// QueryExpression executes a given expression against a given collection and passes the
// result or an error if occured to a given callback handler.
func QueryExpression(exp *Expression, rx ResultCallback) {
	if rx == nil {
		rx = func(err error, _ []bson.M) {}
	}

	wg.Add(1)
	go func() {
		//ensure we decrement the wait counter.
		defer wg.Done()

		// map out the expression condition lists into a bson.M map.
		evalExpr, err := mapExpressionToBSON(exp)

		// if there was an error converting the lists into bson.M maps, then reply
		// the callback and return.
		if err != nil {
			rx(err, nil)
			return
		}

		scl := Session()
		defer scl.Close()

		var result []bson.M

		col := scl.DB(config.Db).C(exp.Collection)

		if err := col.Pipe(evalExpr).All(&result); err != nil {
			rx(err, nil)
			return
		}

		rx(nil, result)
	}()
}

// ResultTypeCallback provides a function type for the returned result from a
// call to Engine.Query, it supplies an error, the expression used in
// and the result of the expression
type ResultTypeCallback func(error, *Expression, []bson.M)

// QueryFile requests a file be loaded which contains specific rules to be used
// in performing a QueryRule,the file contains a json formatted rule set which will
// be processed and results returned to the callback
func QueryFile(file string, params map[string]interface{}, rx ResultTypeCallback) {
	// retrieve the file and load up the content
	qfile, err := os.Open(file)

	// we failed to get file,return err
	if err != nil {
		rx(err, nil, nil)
		return
	}

	// ensure we have the file closed up
	defer qfile.Close()

	rule, err := mapQueryReader2Rule(qfile, params)

	// we failed to create the rule objec,return err
	if err != nil {
		rx(err, nil, nil)
		return
	}

	// execute the given Rule against the concerned database
	QueryRule(rule, rx)
}

// Query takes a byte slice that contain a json rule which gets turned
// rule struct and provided with the needed params to resolve any necessary
// value substitution. End result is supplied to the given callback
func Query(bo []byte, params map[string]interface{}, rx ResultTypeCallback) {
	// generate the rule struct from the given byte slice
	ro, err := mapQuery2Rule(bo, params)

	// if we failed, return and pass to callback
	if err != nil {
		rx(err, nil, nil)
		return
	}

	// execute the given Rule against the concerned database
	QueryRule(ro, rx)
}

// QueryRule runs a given rule set against a collection executing the fail or pass
// depending on the outcome of the Test if no error occured.
func QueryRule(ro *Rule, rx ResultTypeCallback) {
	QueryExpression(ro.Test, func(err error, res []bson.M) {
		if err != nil {
			log.Printf("Error occured executing %s: %s", ro, err)
			return
		}

		if len(res) == 0 {
			QueryExpression(ro.Fail, func(err error, res []bson.M) {
				rx(err, ro.Fail, res)
			})
		} else {
			QueryExpression(ro.Pass, func(err error, res []bson.M) {
				rx(err, ro.Pass, res)
			})
		}
	})
}

//mapQuery2Rule will map out a giving json byte slice into a Rule object.
func mapQuery2Rule(query []byte, params map[string]interface{}) (*Rule, error) {
	return mapQueryReader2Rule(bytes.NewBuffer(query), params)
}

//mapQueryReader2Rule will map out a giving io.Reader into a Rule object.
func mapQueryReader2Rule(query io.Reader, params map[string]interface{}) (*Rule, error) {
	var r Rule

	if err := json.NewDecoder(query).Decode(&r); err != nil {
		return nil, err
	}

	mapAttributesInExpression(r.Test, params)
	mapAttributesInExpression(r.Pass, params)
	mapAttributesInExpression(r.Fail, params)

	return &r, nil
}

// mapAttributesInExpression will take an expression and map into it the giving keys using
// the format '#key#' where found in a supplied map.
func mapAttributesInExpression(exp *Expression, params map[string]interface{}) {
	if params != nil && len(params) != 0 {
		for key, value := range params {
			findKey := fmt.Sprintf("#%s#", key)
			useVal := fmt.Sprintf("%s", value)
			for n, rule := range exp.Conditions {
				exp.Conditions[n] = strings.Replace(rule, findKey, useVal, -1)
			}
		}
	}
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

// mapQuery2Expressions will map out a valid json string through a json decoder
// into a Expression struct. If the param map is supplied, it will substitute the
// appropriate key with the value in the map using the following format '#key#'.
func mapQuery2Expressions(query []byte, params map[string]interface{}) (*Expression, error) {
	var exp = &Expression{}

	if err := json.NewDecoder(bytes.NewBuffer(query)).Decode(exp); err != nil {
		return nil, err
	}

	mapAttributesInExpression(exp, params)
	return exp, nil
}
