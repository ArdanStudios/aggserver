package common

import "gopkg.in/mgo.v2"

// MongoExecuteMux defines a function type, supplied to MongoExecute for performing
// a set action against a giving mongodb session.
type MongoExecuteMux func(c *mgo.Collection) error

// MongoExecute takes a giving mongo session, clones and gets the collection required,
// from the giving database. It calls the giving callback, to execute
// the desired operations, returning an error if such occured.
func MongoExecute(session *mgo.Session, database string, collection string, callback MongoExecuteMux) error {
	// Create a copy of the session, preferable use .Copy(), to ensure to keep the
	// giving previous session's settings.
	newSession := session.Copy()
	defer newSession.Close()

	// Retrieve the database and collection.
	db := newSession.DB(database)
	col := db.C(collection)

	// Executes the giving callback and return the error if non-nil.
	if err := callback(col); err != nil {
		return err
	}

	return nil
}
