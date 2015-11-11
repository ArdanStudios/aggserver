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

// ResultTypeCallback provides a function type for the returned result from a
// call to Engine.Query, it supplies an error, the expression used in
// and the result of the expression
type ResultTypeCallback func(error, *Expression, []bson.M)

// QueryFile requests a file be loaded which contains specific rules to be used
// in performing a QueryRule,the file contains a json formatted rule set which will
// be processed and results returned to the callback
func (e *Engine) QueryFile(file string, params map[string]interface{}, rx ResultTypeCallback) error {
	// retrieve the file and load up the content
	qfile, err := os.Open(file)

	// we failed to get file,return err
	if err != nil {
		// rx(err,nil,nil)
		return err
	}

	// ensure we have the file closed up
	defer qfile.Close()

	rule, err := mapQueryReader2Rule(qfile, params)

	// we failed to create the rule objec,return err
	if err != nil {
		// rx(err,nil,nil)
		return err
	}

	// execute the given Rule against the concerned database
	e.QueryRule(rule, rx)

	return nil
}

// Query takes a byte slice that contain a json rule which gets turned
// rule struct and provided with the needed params to resolve any necessary
// value substitution. End result is supplied to the given callback
func (e *Engine) Query(bo []byte, params map[string]interface{}, rx ResultTypeCallback) error {
	// generate the rule struct from the given byte slice
	ro, err := mapQuery2Rule(bo, params)

	// if we failed, return and pass to callback
	if err != nil {
		// rx(err, nil, nil)
		return err
	}

	// execute the given Rule against the concerned database
	e.QueryRule(ro, rx)
	return nil
}

// QueryRule runs a given rule set against a collection executing the fail or pass
// depending on the outcome of the Test if no error occured.
func (e *Engine) QueryRule(ro *Rule, rx ResultTypeCallback) {
	e.QueryExpression(ro.Test, func(err error, res []bson.M) {
		if err != nil {
			log.Printf("Error occured executing %s: %s", ro, err)
			return
		}

		if len(res) == 0 {
			e.QueryExpression(ro.Fail, func(err error, res []bson.M) {
				rx(err, ro.Fail, res)
			})
		} else {
			e.QueryExpression(ro.Pass, func(err error, res []bson.M) {
				rx(err, ro.Pass, res)
			})
		}
	})
}

// ResultCallback provides a function type for the return of a result from a
// call to Engine.Query.
type ResultCallback func(error, []bson.M)

// QueryExpression executes a given expression against a given collection and passes the
// result or an error if occured to a given callback handler.
func (e *Engine) QueryExpression(exp *Expression, rx ResultCallback) {
	if rx == nil {
		rx = func(err error, _ []bson.M) {}
	}
	go func() {
		//ensure we decrement the wait counter.
		defer e.wg.Done()

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

//mapQuery2Rule will map out a giving json byte slice into a Rule object.
func mapQuery2Rule(query []byte, params map[string]interface{}) (*Rule, error) {
	return mapQueryReader2Rule(bytes.NewBuffer(query), params)
}

//mapQueryReader2Rule will map out a giving io.Reader into a Rule object.
func mapQueryReader2Rule(query io.Reader, params map[string]interface{}) (*Rule, error) {
	r := Rule{
		Test: &Expression{},
		Fail: &Expression{},
		Pass: &Expression{},
	}

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
			useVal := fmt.Sprintf("%q", value)
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
