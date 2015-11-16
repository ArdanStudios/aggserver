package session

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2"
)

// defaultConfig provides a default configuration to be loaded up by Init and
// used to merge into non-provided values in configurations.
var defaultConfig = Config{
	Host:     []string{"ds035428.mongolab.com:35428"},
	AuthDB:   "goinggo",
	Username: "guest",
	Password: "welcome",
	DB:       "goinggo",
}

// Config provides configuration options for the Auth struct
// to create the mongodb session.
type Config struct {
	Host     []string
	Username string
	Password string
	AuthDB   string // Database to use for authentication.
	DB       string // Database to use for session.
}

// mongoSession provides a structure for housing the current mongodb session instance.
type mongoSession struct {
	session *mgo.Session
}

// session provides a global session handler for authentication and entities CRUD
// management.
var session mongoSession

// Init is to be called once and initializes the package session for handling
// model CRUD and authentication requests.
func Init(c *Config) {
	if c == nil {
		c = &defaultConfig
	}

	mgoSession, err := mgo.DialWithInfo(configToDailInfo(c))
	if err != nil {
		panic(err)
	}

	mgoSession.SetMode(mgo.Monotonic, true)

	session.session = mgoSession
}

// Session returns the current active mongodb session.
// Returns a non-nil error if no session was active, i.e
// Init() was not yet called.
func Session() (*mgo.Session, error) {
	if session.session == nil {
		return nil, errors.New("Invalid Session")
	}

	return session.session, nil
}

// MustSession returns the current active mongodb session.
// Panics if Init() has not being called to setup the session yet.
func MustSession() *mgo.Session {
	if session.session == nil {
		panic("Invalid Session, Call Init(*Config)")
	}

	return session.session
}

// configToDailInfo creates a mongo.DialInfo from a given *Config.
func configToDailInfo(c *Config) *mgo.DialInfo {
	return &mgo.DialInfo{
		Addrs:    c.Host,
		Timeout:  60 * time.Second,
		Database: c.AuthDB,
		Username: c.Username,
		Password: c.Password,
	}
}
