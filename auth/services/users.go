package services

import (
	"bytes"
	"encoding/json"

	"github.com/ArdanStudios/aggserver/auth/models"
	"gopkg.in/mgo.v2"
)

// UserEntityService provides a service provider for providing management
// of CRUD and authentication operations for model.UserEntity entities.
type UserEntityService struct{}

// UserService provides a global service providers for UserEntity.
var UserService UserEntityService

// All returns a list of available companies from the underline database.
func (c *UserEntityService) All(session *mgo.Session) ([]*models.UserEntity, error) {
	return models.GetUserEntities(session)
}

// GetUserByEmail gets a particular User using the supplied email data.
func (c *UserEntityService) GetUserByEmail(session *mgo.Session, email string) (*models.UserEntity, error) {
	return models.GetUserByEmail(session, email)
}

// GetUserByPublicID gets a particular User using the supplied serializable data.
// Expects to receive a map with key public_id.
func (c *UserEntityService) GetUserByPublicID(session *mgo.Session, pid string) (*models.UserEntity, error) {
	return models.GetUserByPublicID(session, pid)
}

// Create initializes and creates the the user entity and saves it.
// Expects a map containing properties that match models.UserNew.
func (c *UserEntityService) Create(session *mgo.Session, data []byte) (*models.UserEntity, error) {
	var newUser models.UserNew

	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&newUser)
	if err != nil {
		return err
	}

	user := new(models.UserEntity)
	if err := user.Create(session, &newUser); err != nil {
		return err
	}

	user.SerializeAsPublic()
	return user, nil
}

// ChangePassword resets the giving password with the appropriate crendetials
func (c *UserEntityService) ChangePassword(session *mgo.Session, data []byte) (*models.UserEntity, error) {
	var changeUser models.UserPasswordChange

	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&changeUser)
	if err != nil {
		return err
	}

	user := new(models.UserEntity)
	if err := user.ChangePassword(session, &changeUser); err != nil {
		return err
	}

	user.SerializeAsPublic()
	return user, nil
}

// Destroy destroys the giving entity from the credentails from the map
func (c *UserEntityService) Destroy(session *mgo.Session, data []byte) (*models.UserEntity, error) {
	var destroyUser models.UserDestroy

	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&destroyUser)
	if err != nil {
		return err
	}

	user := new(models.UserEntity)
	if err := user.Destroy(session, &destroyUser); err != nil {
		return err
	}

	user.SerializeAsPublic()
	return user, nil
}

// Update updates the giving entity from the credentails from the map
func (c *UserEntityService) Update(session *mgo.Session, data []byte) (*models.UserEntity, error) {
	var updateUser models.UserUpdate

	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&updateUser)
	if err != nil {
		return err
	}

	user := new(models.UserEntity)
	if err := user.Update(session, &updateUser); err != nil {
		return err
	}

	user.SerializeAsPublic()
	return user, nil
}

// Login authenticates the validity of a entity's credentails from the
// serializable data provided.
func (c *UserEntityService) Login(session *mgo.Session, data []byte) (*models.UserEntity, error) {
	var authUser models.UserLoginAuthentication

	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&authUser)
	if err != nil {
		return err
	}

	user := new(models.UserEntity)
	if err := user.AuthenticateLogin(session, &authUser); err != nil {
		return err
	}

	user.SerializeAsPublic()
	return user, nil
}

// Authenticate authenticates the giving entity from the credentails from the map
func (c *UserEntityService) Authenticate(session *mgo.Session, data []byte) (*models.UserEntity, error) {
	var authUser models.UserAuthentication

	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&authUser)
	if err != nil {
		return err
	}

	user := new(models.UserEntity)
	if err := user.Authenticate(session, &authUser); err != nil {
		return err
	}

	user.SerializeAsPublic()
	return user, nil
}
