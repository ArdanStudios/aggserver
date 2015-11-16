package common

// invalidCrendentailsError is used in the returned error if there was invalid
// crendentails to be validated against.
const InvalidCrendentailsError = "Invalid Authentication Credentials"

// credentailsAuthError is used in the returned error if there was invalid credentail
// data to be authenticated against.
const CredentailsAuthError = "Invalid Authentication Credentials"

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

// UserType is used to indicate type of the user entity.
type UserType int

// constants of the different available company status.
const (
	UnknownType UserType = iota
	BasicUser
	SysAdmin
)

// SerializableEntity defines an interface for serializing an entity into a
// safe structure for public viewaship
type SerializableEntity interface {
	SerializeAsPublic()
}
