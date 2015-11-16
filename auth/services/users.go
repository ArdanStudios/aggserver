package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ArdanStudios/aggserver/auth/common"
	"github.com/ArdanStudios/aggserver/auth/models"
	"github.com/ArdanStudios/aggserver/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UserDatabase defines the database to use.
const UserDatabase = ""

// UserCollection defines the collection to use in storing user entites.
const UserCollection = "users"

// UserPasswordResetCollection defines the collection to use in storing
// user entites password resets.
const UserPasswordResetCollection = "users_password_resets"

// MaxUserLife defines the maximum duration in Hours, of which a reset request remains valid.
const MaxUserResetLife = models.MaxTokenLifeTime * time.Hour

// userService provides a service provider for providing management
// of CRUD and authentication operations for model.User entities.
type userService struct{}

// UserService provides a global service providers for User entity management.
var UserService userService

// AllUsers returns a list of available companies from the underline database.
// Returns the users lists if found or else returns a non-nil error.
func (c *userService) AllUsers(session *mgo.Session) ([]*models.User, error) {
	log.Dev("userService", "userService.AllUsers", "Started")
	var entities []*models.User

	log.User("userService", "userService.AllUsers", "Load all users")
	err := common.MongoExecute(session, UserDatabase, UserCollection, func(co *mgo.Collection) error {
		log.Dev("userService", "userService.AllUsers", "Completed")
		return co.Find(nil).All(&entity)
	})

	if err != nil {
		log.Dev("userService", "userService.AllUsers", "Database %s : Collection %s : Error %s", UserDatabase, UserCollection, err.Error())
		return nil, err
	}

	return entities, nil
}

// GetUserByEmail gets a particular User using the supplied email data.
// Returns the user if found or else returns a non-nil error.
func (c *userService) GetUserByEmail(session *mgo.Session, email string) (*models.User, error) {
	log.Dev("userService", "userService.GetUserByEmail", "Started For User : Email %s", email)
	var entity models.User

	log.User("userService", "userService.GetUserByEmail", "Retrieve User with Email Email %s", email)
	if err := common.MongoExecute(session, UserDatabase, UserCollection, func(co *mgo.Collection) error {
		log.Dev("userService", "userService.GetUserByEmail", "Completed For User : Email %s", email)
		return co.Find(bson.M{"email": email}).One(&entity)
	}); err != nil {
		log.Dev("userService", "userService.GetUserByEmail", "Started For User : Email %s : Error %s", email, err.Error())
		return nil, err
	}

	return &entity, nil
}

// GetUserByPublicID gets a particular User using the supplied PublicID of the entity.
// Returns the user if found or else returns a non-nil error.
func (c *userService) GetUserByPublicID(session *mgo.Session, pid string) (*models.User, error) {
	log.Dev("userService", "userService.GetUserByPublicID", "Started For User ID %s", pid)
	var entity models.User

	log.User("userService", "userService.GetUserByPublicID", "User ID %s", pid)
	if err := common.MongoExecute(session, UserDatabase, UserCollection, func(co *mgo.Collection) error {
		log.Dev("userService", "userService.GetUserByPublicID", "Completed For User ID %s", pid)
		return co.Find(bson.M{"public_id": publicID}).One(&entity)
	}); err != nil {
		log.Dev("userService", "userService.GetUserByPublicID", "User ID %s : Error : %s", err.Error())
		return nil, err
	}

	return &entity, nil
}

// Save inserts a given user entity record from the provided serilizable data.
// Returns a non-nil error if the operation fails, else returns the generated user
// entity.
func (c *userService) Save(session *mgo.Session, data []byte) (*models.User, error) {
	log.Dev("userService", "userService.Save", "Started")
	var user models.User

	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&user)
	if err != nil {
		log.Dev("userService", "userService.Save", "Session : Create : Error : %s", err.Error())
		return nil, err
	}

	meta := fmt.Sprintf("{FirstName %q | Last Name %s | Email %s}", user.FirstName, user.LastName, user.Email)

	log.User("userService", "userService.Save", "Session : Save : User.SetToken : Email %s", user.Email)
	if err := user.SetToken(); err != nil {
		log.Dev("userService", "userService.Save", "Session : User.SetToken : Email %s: Error : %s", user.Email, err.Error())
		return nil, err
	}

	if err := common.MongoExecute(session, UserDatabase, UserCollection, func(co *mgo.Collection) error {
		log.Dev("userService", "userService.Save", "Finished for User : Meta %s : Email %s", meta, u.Email)
		return co.Insert(&user)
	}); err != nil {
		log.Dev("userService", "userService.Save", "Save User : Meta %s : Error : %s", meta, err.Error())
		return nil, err
	}

	return nil, err
}

// Create initializes and saves the the user entity and saves it.
// Expects a serializable data containing the needed credentials,
// that match a models.UserNew. Returns a non-nil error if the operation fails,
// else returns the user entity.
func (c *userService) Create(session *mgo.Session, data []byte) (*models.User, error) {
	log.Dev("userService", "userService.Create", "Started")
	var newUser models.UserNew

	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&newUser)
	if err != nil {
		log.Dev("userService", "userService.Create", "Session : Create : Error : %s", err.Error())
		return nil, err
	}

	meta := fmt.Sprintf("{ Name: %q  Email: %q }", fmt.Sprintf("%s %s", newUser.FirstName, newUser.LastName), newUser.Email)
	log.User("userService", "userService.Create", "Service : Entity : Create : %s", meta)

	user := new(models.User)
	if err := user.Create(session, &newUser); err != nil {
		log.Dev("userService", "userService.Create", "Session : Create %s : Error : %s", meta, err.Error())
		return nil, err
	}

	// Insert this record into the database.
	if err := common.MongoExecute(session, UserDatabase, UserCollection, func(co *mgo.Collection) error {
		log.Dev("userService", "userService.Create", "Completed")
		return co.Insert(c)
	}); err != nil {
		log.Dev("userService", "userService.Create", "Mongodb.Insert() : Create %s : Error : %s", meta, err.Error())
		return nil, err
	}

	user.SerializeAsPublic()
	return user, nil
}

// ForgetPassword initiates a password recovery process, by registering a password
// reset instruction for a specific duration, until either that duration expires or a password
// change requests is recieved for the given user.
// Returns a non-nil error if the operation fails, else returns the user entity.
func (c *userService) ForgetPassword(session *mgo.Session, data []byte) (*models.User, error) {
	log.Dev("userService", "userService.ForgetPassword", "Started")
	var forgetUser models.UserPasswordReset

	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&forgetUser)
	if err != nil {
		return nil, err
	}

	log.Dev(forgetUser.PublicID, "userService.ForgetPassword", "Load User Started")
	var user models.User

	// Check if we indeed have the correct details for a user.
	if err := common.MongoExecute(session, UserDatabase, UserCollection, func(co *mgo.Collection) error {
		log.Dev(forgetUser.PublicID, "userService.ForgetPassword", "Load User Completed")
		return co.Find(bson.M{"public_id": forgetUser.PublicID}).One(&user)
	}); err != nil {
		return nil, err
	}

	meta := fmt.Sprintf("{ Name: %q  Email: %q }", fmt.Sprintf("%s %s", user.FirstName, user.LastName), user.Email)
	log.User(forgetUser.PublicID, "userService.ForgetPassword", "Service : Entity : ForgetPassword : %s", meta)

	// Check whether we have already a reset request pending.
	var pendingReset models.UserPasswordReset

	log.Dev(forgetUser.PublicID, "userService.ForgetPassword", "Pending : UserPasswordReset")

	if err := common.MongoExecute(session, UserDatabase, UserPasswordResetCollection, func(co *mgo.Collection) error {
		log.Dev(forgetUser.PublicID, "userService.ForgetPassword", "Pending : UserPasswordReset : Completed")
		return co.Find(bson.M{"public_id": resetUser.PublicID, "token": resetUser.Token}).One(&pendingReset)
	}); err != nil {
		log.Dev(forgetUser.PublicID, "userService.ForgetPassword", "Pending : UserPasswordReset : Error  %s", err.Error())
		return nil, err
	}

	log.User(forgetUser.PublicID, "userService.ForgetPassword", "Pending : Check Expiration")

	// If there is one pending, then check if the time has not expired.
	ms := time.Since(pendingReset.ExpireAt)
	log.Dev(forgetUser.PublicID, "userService.ForgetPassword", "Pending : Check Expiration : Expired at %q : Elapsed %q", pendingReset.ExpireAt, ms)

	if ms.Hours() > MaxUserResetLife {
		// The pending request as expired. Remove.
		log.User(pendingReset.PublicID, "userService.ForgetPassword", "Started : Pending : Remove Expired Reset")
		if err := common.MongoExecute(session, UserDatabase, UserPasswordResetCollection, func(co *mgo.Collection) error {
			log.Dev(pendingReset.PublicID, "userService.ForgetPassword", "Completed : Pending : Remove Expired Reset")
			return co.RemoveId(pendingReset.ID)
		}); err != nil {
			log.Dev(pendingReset.PublicID, "userService.ForgetPassword", "Error %s : Pending : Remove Expired Reset", err.Error())
			return nil, err
		}
	} else {
		log.User(forgetUser.PublicID, "userService.ForgetPassword", "Found Pending UserPasswordReset")
		return nil, errors.New("Pending Password Reset")
	}

	log.User(forgetUser.PublicID, "userService.ForgetPassword", "Adding UserPasswordReset ")

	// Add the latest reset requests
	if err := common.MongoExecute(session, UserDatabase, UserPasswordResetCollection, func(co *mgo.Collection) error {
		log.Dev(forgetUser.PublicID, "userService.ForgetPassword", "Completed")
		return co.Insert(&forgetUser)
	}); err != nil {
		log.Dev(forgetUser.PublicID, "userService.ForgetPassword", "Error %s", err.Error())
		return nil, err
	}

	user.SanitizeAsPublic()
	return &user, nil
}

// ResetPassword initiates a password reset on the giving entity by using it's
// PublicID. It checks if any previous pending password requests exists and
// fullfills the requests if the requests as not expired.
// Returns a non-nil error if the operation fails, else returns the user entity.
func (c *userService) ResetPassword(session *mgo.Session, data []byte) (*models.User, error) {
	log.Dev("userService", "userService.ResetPassword", "Started")
	var changeUser models.UserPasswordChange

	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&changeUser)
	if err != nil {
		log.Dev("userService", "userService.ResetPassword", "JSON Decode : Error %s", err.Error())
		return nil, err
	}

	if err := changeUser.ValidatePassword(); err != nil {
		log.Dev(changeUser.PublicID, "userService.ResetPassword", "Error : UserPasswordChange.Validate : %s", err.Error())
		return nil, err
	}

	var pending UserPasswordReset

	log.User(changeUser.PublicID, "userService.ResetPassword", "Load UserPasswordReset")
	if err := common.MongoExecute(session, UserDatabase, UserPasswordResetCollection, func(co *mgo.Collection) error {
		log.Dev(changeUser.PublicID, "userService.ResetPassword", "Load UserPasswordReset : Complete")
		return co.Find(bson.M{"public_id": changeUser.PublicID, "token": changeUser.Token}).One(&pendingReset)
	}); err != nil {
		log.Dev(changeUser.PublicID, "userService.ResetPassword", "Load UserPasswordReset : Error %s", err.Error())
		return nil, err
	}

	ms := time.Since(pending.ExpireAt)
	log.User(changeUser.PublicID, "userService.ResetPassword", "PasswordReset : Time %s : Elapsed %s", pending.ExpiredAt, ms)
	if ms.Hours() > MaxUserResetLife {
		log.Dev(changeUser.PublicID, "userService.ResetPassword", "PasswordReset : Expired")
		return nil, errors.New("PasswordReset Expired")
	}

	var user models.User
	log.Dev(changeUser.PublicID, "userService.ResetPassword", "User : Mongodb.Find() : Started")
	if err := common.MongoExecute(session, UserDatabase, UserCollection, func(co *mgo.Collection) error {
		log.Dev(changeUser.PublicID, "userService.ResetPassword", "User : Mongodb.Find() : Completed")
		return co.Find(bson.M{"public_id": changeUser.PublicID}).One(&user)
	}); err != nil {
		log.Dev(changeUser.PublicID, "userService.ResetPassword", "User : Mongodb.Find() : Error %s", err.Error())
		return nil, err
	}

	log.User(changeUser.PublicID, "userService.ResetPassword", "User : ChangePassword")
	if err := user.ChangePassword(session, &changeUser); err != nil {
		log.Dev(changeUser.PublicID, "userService.ResetPassword", "User : ChangePassword : Error %s", err.Error())
		return nil, err
	}

	q := bson.M{"public_id": c.PublicID}
	mod := bson.M{"password": c.Password, "modified_at": c.ModifiedAt}

	log.Dev(changeUser.PublicID, "userService.ResetPassword", "User : Mongodb.Update() : Started")
	if err := common.MongoExecute(session, UserDatabase, UserCollection, func(co *mgo.Collection) error {
		log.Dev(changeUser.PublicID, "userService.ResetPassword", "User : Mongodb.Update() : Completed")
		return co.Update(q, bson.M{"$set": mod})
	}); err != nil {
		return nil, err
	}

	// Get the required expired time stamp and clean out all reset requests
	// within the range.
	now := time.Now()
	expired := now.Add(MaxUserResetLife)
	q := bson.M{"expired_at": bson.M{"$lt": expired}}

	log.Dev(changeUser.PublicID, "userService.ResetPassword", "UserPasswordReset Collection : Mongodb.RemovalAll() : Started")
	// Remove all pending requests regardless of ID that has expired.
	if err := common.MongoExecute(session, UserDatabase, UserPasswordResetCollection, func(co *mgo.Collection) error {
		log.Dev(changeUser.PublicID, "userService.ResetPassword", "UserPasswordReset Collection : Mongodb.RemovalAll() : Completed")
		return co.RemoveAll(q)
	}); err != nil {
		log.Dev(changeUser.PublicID, "userService.ResetPassword", "UserPasswordReset Collection : Mongodb.RemovalAll() : Error %s", err.Error())
		return nil, err
	}

	user.SerializeAsPublic()
	log.Dev("userService", "userService.ResetPassword", "Completed")
	return user, nil
}

// ChangePassword resets the giving password with the appropriate credentails from
// the supplied serializable data.
// Returns a non-nil error if the operation fails, else returns the user entity.
func (c *userService) ChangePassword(session *mgo.Session, data []byte) (*models.User, error) {
	log.Dev("userService", "userService.ChangePassword", "Started")
	var changeUser models.UserPasswordChange

	log.Dev("userService", "userService.ChangePassword", "JSON Decode : Started")
	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&changeUser)
	if err != nil {
		log.Dev("userService", "userService.ChangePassword", "JSON Decode : Error %s", err.Error())
		return nil, err
	}

	log.User(changeUser.PublicID, "userService.ChangePassword", "User : ValidatePassword")
	if err := changeUser.ValidatePassword(); err != nil {
		log.Dev(changeUser.PublicID, "userService.ChangePassword", "User : ValidatePassword : Error %s", err.Error())
		return nil, err
	}

	var user models.User

	log.Dev(changeUser.PublicID, "userService.ChangePassword", "Started : User : UserDatabase %s: UserCollection %s : Mongodb.Find().One()", UserDatabase, UserCollection)
	if err := common.MongoExecute(session, UserDatabase, UserCollection, func(co *mgo.Collection) error {
		log.Dev(changeUser.PublicID, "userService.ChangePassword", "Completed : User : UserDatabase %s: UserCollection %s : Mongodb.Find().One()", UserDatabase, UserCollection)
		return co.Find(bson.M{"public_id": changeUser.PublicID}).One(&user)
	}); err != nil {
		log.Dev(changeUser.PublicID, "userService.ChangePassword", "User : UserDatabase %s: UserCollection %s : Mongodb.Find().One() : Error %s", UserDatabase, UserCollection, err.Error())
		return nil, err
	}

	log.User(changeUser.PublicID, "userService.ChangePassword", "User : User.ChangePassword")
	if err := user.ChangePassword(session, &changeUser); err != nil {
		log.Dev(changeUser.PublicID, "userService.ChangePassword", "User : User.ChangePassword : Error %s", err.Error())
		return nil, err
	}

	q := bson.M{"public_id": c.PublicID}
	mod := bson.M{"password": c.Password, "modified_at": c.ModifiedAt}

	if err := common.MongoExecute(session, UserDatabase, UserCollection, func(co *mgo.Collection) error {
		log.Dev(user.PublicID, "userService.ChangePassword", "Completed")
		return co.Update(q, bson.M{"$set": mod})
	}); err != nil {
		log.Dev(user.PublicID, "userService.ChangePassword", "Error %s", err.Error())
		return nil, err
	}

	user.SerializeAsPublic()
	return user, nil
}

// Destroy destroys the giving entity from the credentails from the serializable
// data. It uses the entity's PublicID.
// Returns a non-nil error if the operation fails, else returns the user entity.
func (c *userService) Destroy(session *mgo.Session, data []byte) (*models.User, error) {
	log.Dev("userService", "userService.Destroy", "Started")
	var destroyUser models.UserDestroy

	log.Dev("userService", "userService.Destroy", "JSON.Decode")
	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&destroyUser)
	if err != nil {
		log.Dev("userService", "userService.Destroy", "JSON.Decode : Error %s", err.Error())
		return nil, err
	}

	var user models.User

	log.Dev(destroyUser.PublicID, "userService.Destroy", "User : Load User : Mongodb.Find().One()")
	if err := common.MongoExecute(session, UserDatabase, UserCollection, func(co *mgo.Collection) error {
		log.Dev(destroyUser.PublicID, "userService.Destroy", "Completed : User : Load User : Mongodb.Find().One()")
		return co.Find(bson.M{"public_id": changeUser.PublicID}).One(&user)
	}); err != nil {
		log.Dev(destroyUser.PublicID, "userService.Destroy", "User : Load User : Mongodb.Find().One() : Error %s", err.Error())
		return nil, err
	}

	if err := common.MongoExecute(session, UserDatabase, UserCollection, func(co *mgo.Collection) error {
		log.Dev("userService", "userService.Destroy", "Completed")
		return co.RemoveId(c.ID)
	}); err != nil {
		log.Dev("userService", "userService.Destroy", "Error %s", err.Error())
		return err
	}

	user.SerializeAsPublic()
	return user, nil
}

// Update updates the giving entity from the credentails from the serializable
// data.
// Returns a non-nil error if the operation fails, else returns the user entity.
func (c *userService) Update(session *mgo.Session, data []byte) (*models.User, error) {
	log.Dev("userService", "userService.Update", "Started")
	var updateUser models.UserUpdate

	log.Dev("userService", "userService.Update", "JSONDecoder.Decode ")
	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&updateUser)
	if err != nil {
		log.Dev("userService", "userService.Update", "JSONDecoder.Decode : Error %s", err.Error())
		return nil, err
	}

	var user models.User

	log.Dev(updateUser.PublicID, "userService.Update", "User : UserDatabase %s : UserCollection %s : Mongodb.Find().One()", UserDatabase, UserCollection)
	if err := common.MongoExecute(session, UserDatabase, UserCollection, func(co *mgo.Collection) error {
		log.Dev(updateUser.PublicID, "userService.Update", "Complete : User : UserDatabase %s : UserCollection %s : Mongodb.Find().One()", UserDatabase, UserCollection)
		return co.Find(bson.M{"public_id": changeUser.PublicID}).One(&user)
	}); err != nil {
		log.Dev(updateUser.PublicID, "userService.Update", "Complete Error %s : User : UserDatabase %s : UserCollection %s : Mongodb.Find().One()", UserDatabase, UserCollection, err.Error())
		return nil, err
	}

	log.User(updateUser.PublicID, "userService.Update", "User : User.Update()")
	if err := user.Update(session, &updateUser); err != nil {
		log.Dev(updateUser.PublicID, "userService.Update", "User : User.Update() : Error %s", err.Error())
		return nil, err
	}

	if err := common.MongoExecute(session, UserDatabase, UserCollection, func(co *mgo.Collection) error {
		log.Dev("userService", "userService.Update", "Completed")
		return co.UpdateId(c.ID, c)
	}); err != nil {
		log.Dev("userService", "userService.Update", "Completed Error %s", err.Error())
		return nil, err
	}

	user.SerializeAsPublic()

	return user, nil
}

// Login authenticates the validity of a entity's credentails from the
// serializable data provided. It uses the entity's Email and Password for
// authentication.
// Returns a non-nil error if the authentication failed, else returns the user entity.
func (c *userService) Login(session *mgo.Session, data []byte) (*models.User, error) {
	log.Dev("userService", "userService.Login", "Started")
	var authUser models.UserLoginAuthentication

	log.Dev("userService", "userService.Login", "User : JSONDecoder.Decode")
	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&authUser)
	if err != nil {
		log.Dev("userService", "userService.Login", "User : JSONDecoder.Decode : Error %s", err.Error())
		return nil, err
	}

	var user models.User

	log.User("userService", "userService.Login", "User : LoadUser : Email %s", authUser.Email)
	if err := common.MongoExecute(session, UserDatabase, UserCollection, func(co *mgo.Collection) error {
		log.Dev(authUser.Email, "userService.Login", "User : Mongodb.Find().One : Success")
		return co.Find(bson.M{"public_id": changeUser.PublicID}).One(&user)
	}); err != nil {
		log.Dev(authUser.Email, "userService.Login", "User : Mongodb.Find().One : Error %s", err.Error())
		return nil, err
	}

	log.User(authUser.Email, "userService.Login", "AuthenticateLogin: Started")
	if err := user.AuthenticateLogin(session, &authUser); err != nil {
		log.Dev(authUser.Email, "userService.Login", "AuthenticateLogin : Completed Error %s", err.Error())
		return nil, err
	}

	log.Dev("userService", "userService.Login", "Completed")
	user.SerializeAsPublic()
	return user, nil
}

// Authenticate authenticates the giving entity from the credentails from the
// serializable data. It uses the entity's PublicID and Token for authentication.
// Returns a non-nil error if the authentication failed, else returns the user entity.
func (c *userService) Authenticate(session *mgo.Session, data []byte) (*models.User, error) {
	log.Dev("userService", "userService.Authenticate", "Started")
	var authUser models.UserAuthentication

	log.Dev("userService", "userService.Authenticate", "User : JSONDecoder.Decode")
	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&authUser)
	if err != nil {
		log.Dev("userService", "userService.Authenticate", "User : JSONDecoder.Decode : Error %s", err.Error())
		return nil, err
	}

	var user models.User

	log.User("userService", "userService.Authenticate", "User : LoadUser : Email %s", authUser.Email)
	if err := common.MongoExecute(session, UserDatabase, UserCollection, func(co *mgo.Collection) error {
		log.Dev(authUser.Email, "userService.Authenticate", "User : Mongodb.Find().One : Success")
		return co.Find(bson.M{"public_id": changeUser.PublicID}).One(&user)
	}); err != nil {
		log.Dev(authUser.Email, "userService.Authenticate", "AuthenticateLogin : Failed : Error %s", err.Error())
		return nil, err
	}

	log.User(authUser.Email, "userService.Authenticate", "AuthenticateToken: Started")
	if err := user.AuthenticateToken(session, &authUser); err != nil {
		log.Dev(authUser.Email, "userService.Authenticate", "AuthenticateToken : Completed Error %s", err.Error())
		return nil, err
	}

	log.Dev("userService", "userService.Authenticate", "Completed")
	user.SerializeAsPublic()
	return user, nil
}
