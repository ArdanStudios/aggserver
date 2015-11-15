package auth

import (
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

// mgoSession provides the global mongodb session for handling auth requests.
var mgoSession mgo.Session

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
