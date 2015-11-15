package auth

import (
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

// Config provides configuration options for the Auth struct
// to create the mongodb session.
type Config struct {
	Host     []string
	Username string
	Password string
	AuthDB   string // Database to use for authentication.
	DB       string // Database to use for session.
}

// ConfigToDailInfo creates a mongo.DialInfo from a given
// *Config.
func ConfigToDailInfo(c *Config) *mgo.DialInfo {
	return &mgo.DialInfo{
		Addrs:    c.Host,
		Timeout:  60 * time.Second,
		Database: c.AuthDB,
		Username: c.Username,
		Password: c.Password,
	}
}

// Service provides an interface for a authentication service handler to
// handle requests for a specific service.
type Service interface {
	Create(*mgo.Session, []byte) error
	Update(*mgo.Session, []byte) error
	Destroy(*mgo.Session, []byte) error
	Authenticate(*mgo.Session, []byte) error
}

// Auth connects to a mongodb session and handles managment of entities CRUD and
// authentication using service providers that meet the Service interface.
type Auth struct {
	*Config
	*mgo.Session
	serviceMutex     sync.RWMutex
	serviceProviders map[string]Service // Map of entity service providers.
}
