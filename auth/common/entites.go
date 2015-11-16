package common

import "gopkg.in/mgo.v2"

// invalidCrendentailsError is used in the returned error if there was invalid
// crendentails to be validated against.
const invalidCrendentailsError = "Invalid Authentication Credentials"

// credentailsAuthError is used in the returned error if there was invalid credentail
// data to be authenticated against.
const credentailsAuthError = "Invalid Authentication Credentials"

// EntityStatus is used to indicate current status of the entity.
type EntityStatus int

// constants of the different available company status.
const (
	NoStatusEntity EntityStatus = iota
	InactiveEntity
	ActiveEntity
	DisabledEntity
	DestroyedEntity
)

// EntityCreate defines an interface for creating a new entity from a serializable data
// and inserting it into a collection using a mongodb session.
type EntityCreate interface {
	Create(*mgo.Session, []byte) error
}

// EntityInsert defines an interface for inserting a new entity into a mongodb
// collection using a supplied mongodb session.
type EntityInsert interface {
	Insert(*mgo.Session) error
}

// EntityDestroy defines an interface for the removal/destruction of an entity
// using a mongo session.
type EntityDestroy interface {
	Destroy(*mgo.Session) error
}

// EntityUpdate defines an interface for the updating an entity
// using a mongo session and its serializable update data.
type EntityUpdate interface {
	Update(*mgo.Session, []byte) error
}

// EntityAuthenticate defines an interface for the updating an entity
// using a mongo session and its serializable update data.
type EntityAuthenticate interface {
	Authenticate(*mgo.Session, []byte) error
}

// EntityAuthable defines a composed interface of all operations an entity
// capable of creation, updating, destruction and authentication.
type EntityAuthable interface {
	EntityCreate
	EntityInsert
	EntityUpdate
	EntityDestroy
	EntityAuthenticate
}

// SerializableEntity defines an interface for serializing an entity into a
// safe structure for public viewaship
type SerializableEntity interface {
	SerializeAsPublic()
}
