package models

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"
	"unicode"

	"github.com/ArdanStudios/aggserver/auth/common"
	"github.com/ArdanStudios/aggserver/auth/vendor/github/satori/go.uuid"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//InvalidPasswordError is returned when the password given fails to match a given
// set of requirments.
const InvalidPasswordError = "Invalid Password"

// MinPasswordLength defines the minimum length a password should have.
const MinPasswordLength = 8

// MaxPasswordLength defines the maximum length a password should have.
const MaxPasswordLength = 50

// MaxTokenLifeTime defines the maximum time a user's token is considered valid.
// Its value is giving in hours.
const MaxTokenLifeTime = 24

// UserEntityDatabase defines the database to use.
const UserEntityDatabase = ""

// UserEntityCollection defines the collection to use in storing user entites.
const UserEntityCollection = "users"

// UserEntityPasswordResetCollection defines the collection to use in storing
// user entites password resets.
const UserEntityPasswordResetCollection = "users_password_resets"

// UserEntity represents an entity for user records.
type UserEntity struct {
	ID         bson.ObjectId `json:"id" bson:"id"`
	FirstName  string        `json:"first_name,omitempty" bson:"first_name"`
	LastName   string        `json:"last_name,omitempty" bson:"last_name"`
	Token      string        `json:"token,omitempty" bson:"-"`
	PublicID   string        `json:"public_id,omitempty" bson:"public_id"`
	PrivateID  string        `json:"private_id,omitempty" bson:"private_id"`
	Email      string        `bson:"email" json:"email"`
	Password   string        `bson:"password" json:"password,omitempty"`
	Status     EntityStatus  `json:"status" bson:"status"`
	ModifiedAt *time.Time    `json:"modified_at,omitempty" bson:"modified_at"`
	CreatedAt  *time.Time    `json:"created_at,omitempty" bson:"created_at"`
}

// Insert inserts the Company entity into the mongoDB database collection.
// It returns a non-nil error if an error occurs
func (c *UserEntity) Insert(session *mgo.Session) error {
	return common.MongoExecute(session, UserEntityDatabase, UserEntityCollection, func(co *mgo.Collection) error {
		return co.Insert(c)
	})
}

// UserNew provides a struct for use in creating a new company entity.
type UserNew struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

// ValidatePassword validates the new password and password confirmation are a
// match and that they follow the given requirements for a valid password.
func (u *UserNew) ValidatePassword() error {
	if err := ValidatePassword(u.Password); err != nil {
		return err
	}

	if u.Password != u.PasswordConfirm {
		return errors.New(InvalidPasswordError)
	}

	return nil
}

// Create defines a new company entity and saves it into the giving entity database
// using the provided mongo session and serializable data.
func (c *UserEntity) Create(session *mgo.Session, newUser *UserNew) error {
	publicUUID := uuid.NewV4()
	privateUUID := uuid.NewV4()

	createdAt := time.Now()
	modifiedAt := time.Now()

	c.FirstName = newUser.FirstName
	c.LastName = newUser.LastName
	c.Email = newUser.Email
	c.PublicID = publicUUID.String()
	c.PrivateID = privateUUID.String()
	c.CreatedAt = createdAt
	c.ModifiedAt = modifiedAt
	c.Status = EntityActive

	// create the user password hash
	p, err := crypto.BcryptHash(c.PrivateID + newUser.Password)
	if err != nil {
		return err
	}

	c.Password = string(p)

	// Set up the entity's authentication token.
	if err := c.SetToken(); err != nil {
		return nil, err
	}

	// Insert this record into the database.
	return common.MongoExecute(session, UserEntityDatabase, UserEntityCollection, func(co *mgo.Collection) error {
		return co.Insert(c)
	})
}

// UserPasswordReset is used to resets a users entity's password
type UserPasswordReset struct {
	ID       bson.ObjectId `bson:"id" json:"id,omitempty"`
	PublicID string        `bson:"public_id" json:"public_id"`
	Token    string        `bson:"token" json:"token"`
	ExpireAt *time.Time    `bson:"expire_at" json:"expire_at"`
}

// ForgetPassword registers a password forget instruction for the given entity.
func (c *UserEntity) ForgetPassword(session *mgo.Session, resetUser *UserPasswordReset) error {
	// Check if credentails are not loaded then load the document.
	if !c.isCredentailsLoaded() {
		if err := common.MongoExecute(session, UserEntityDatabase, UserEntityCollection, func(co *mgo.Collection) error {
			return co.Find(bson.M{"public_id": resetUser.PublicID}).One(c)
		}); err != nil {
			return err
		}
	}

	if c.PublicID != resetUser.PublicID {
		return errors.New(invalidCrendentailsError)
	}

	// Check if we have already a reset request pending.
	var pendingReset UserPasswordReset
	if err := common.MongoExecute(session, UserEntityDatabase, UserEntityPasswordResetCollection, func(co *mgo.Collection) error {
		return co.Find(bson.M{"public_id": resetUser.PublicID, "token": resetUser.Token}).One(&pendingReset)
	}); err != nil {
		return err
	}

	// If there is a pending then check if the time has not expired.
	if pendingReset != nil {
		if time.Since(pendingReset.ExpireAt) > (MaxTokenLifeTime * time.Hour) {
			// We are passed the given time, so expire
			if err := common.MongoExecute(session, UserEntityDatabase, UserEntityPasswordResetCollection, func(co *mgo.Collection) error {
				return co.RemoveId(pendingReset.ID)
			}); err != nil {
				return err
			}
		} else {
			return errors.New("Pending Password Reset")
		}
	}

	// Add the latest reset requests
	if err := common.MongoExecute(session, UserEntityDatabase, UserEntityPasswordResetCollection, func(co *mgo.Collection) error {
		return co.Insert(resetUser)
	}); err != nil {
		return err
	}

	c.SerializeAsPublic()
	return nil
}

// UserPasswordChange is used to provide password change instruction for a user
// entity.
type UserPasswordChange struct {
	Token           string `json:"token"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"confirm_password"`
	PublicID        string `json:"public_id"`
}

// ResetPassword resets the password according to the values being received.
// It also clears all password pending reset requests.
func (c *UserEntity) ResetPassword(session *mgo.Session, changeUser *UserPasswordChange) error {
	var pending UserPasswordReset

	if err := common.MongoExecute(session, UserEntityDatabase, UserEntityPasswordResetCollection, func(co *mgo.Collection) error {
		return co.Find(bson.M{"public_id": changeUser.PublicID, "token": changeUser.Token}).One(&pendingReset)
	}); err != nil {
		return err
	}

	// Perform the necessary password reset
	if err := c.ChangePassword(session, changeUser); err != nil {
		return err
	}

	// Get the required expired time stamp and clean out all reset requests
	// within the range.
	now := time.Now()
	expired := now.Add(time.Hour * -MaxTokenLifeTime)
	q := bson.M{"expired_at": bson.M{"$lt": expired}}

	// Remove all pending requests regardless of ID that has expired.
	if err := common.MongoExecute(session, UserEntityDatabase, UserEntityPasswordResetCollection, func(co *mgo.Collection) error {
		return co.RemoveAll(q)
	}); err != nil {
		return err
	}

	return nil
}

// ChangePassword changes the password according to the values being received.
func (c *UserEntity) ChangePassword(session *mgo.Session, changeUser *UserPasswordChange) error {
	if err := changeUser.ValidatePassword(); err != nil {
		return err
	}

	// Check if credentails are not loaded then load the document.
	if !c.isCredentailsLoaded() {
		if err := common.MongoExecute(session, UserEntityDatabase, UserEntityCollection, func(co *mgo.Collection) error {
			return co.Find(bson.M{"public_id": changeUser.PublicID}).One(c)
		}); err != nil {
			return err
		}
	}

	if resetUser.PublicID != c.PublicID {
		return errors.New(invalidCrendentailsError)
	}

	p, err := crypto.BcryptHash(c.PrivateID + changeUser.Password)
	if err != nil {
		return err
	}

	c.Password = string(p)
	c.ModifiedAt = time.Now()

	q := bson.M{"public_id": c.PublicID}
	mod := bson.M{"password": c.Password, "modified_at": c.ModifiedAt}

	if err := common.MongoExecute(session, UserEntityDatabase, UserEntityCollection, func(co *mgo.Collection) error {
		return co.Update(q, bson.M{"$set": mod})
	}); err != nil {
		return err
	}

	if err := c.SetToken(); err != nil {
		return err
	}

	// c.SerializeAsPublic()
	return nil
}

// ValidatePassword validates the new password and password confirmation are a
// match and that they follow the given requirements for a valid password.
func (u *UserPasswordChange) ValidatePassword() error {
	if err := ValidatePassword(u.Password); err != nil {
		return err
	}

	if u.Password != u.PasswordConfirm {
		return errors.New(InvalidPasswordError)
	}

	return nil
}

// UserUpdate provides a struct for use updating an existing company entity.
type UserUpdate struct {
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	Email     string       `json:"email"`
	Username  string       `json:"username"`
	Password  string       `json:"password"`
	PublicID  string       `json:"public_id"`
	Status    EntityStatus `json:"status" `
	Token     string       `json:"token"`
}

// Update defines an update to a company's entity and updates the giving entity
// in the corresponding mongodb database and giving collection using the provided
// mongo session and serializable data.
func (c *UserEntity) Update(session *mgo.Session, updatingUser *UserUpdate) error {
	// Check if credentails are not loaded then load the document.
	if !c.isCredentailsLoaded() {
		if err := common.MongoExecute(session, UserEntityDatabase, UserEntityCollection, func(co *mgo.Collection) error {
			return co.Find(bson.M{"public_id": updatingUser.PublicID}).One(c)
		}); err != nil {
			return err
		}
	}

	if updatingUser.PublicID != c.PublicID {
		return errors.New("Invalid PublicID for entity")
	}

	c.Name = updatingUser.Name
	c.Status = updatingUser.Status

	if updatingUser.Config != nil {
		c.Config = updatingUser.Config
	}

	c.ModifiedAt = time.Now()

	return common.MongoExecute(session, UserEntityDatabase, UserEntityCollection, func(co *mgo.Collection) error {
		return co.UpdateId(c.ID, c)
	})
}

// UserDestroy contains the given user public_id to instruct a removal of the entity
// from the databse.
type UserDestroy struct {
	PublicID string `json:"public_id"`
}

// Destroy destroys/removes a company entity from the corresponding mongodb database
// and giving collection using the provided mongo session and serializable data.
func (c *UserEntity) Destroy(session *mgo.Session, destroyUser *UserDestroy) error {
	// Check if credentails are not loaded then load the document.
	if !c.isCredentailsLoaded() {
		if err := common.MongoExecute(session, UserEntityDatabase, UserEntityCollection, func(co *mgo.Collection) error {
			return co.Find(bson.M{"public_id": destroyUser.PublicID}).One(c)
		}); err != nil {
			return err
		}
	}

	if err := common.MongoExecute(session, UserEntityDatabase, UserEntityCollection, func(co *mgo.Collection) error {
		return co.RemoveId(c.ID)
	}); err != nil {
		return err
	}

	return nil
}

// UserAuthentication provides the necessary credentials needed for authenticating
// a user entity.
type UserAuthentication struct {
	PublicID string `json:"public_id"`
	Token    string `json:"public_id"`
}

// Authenticate authenticates the token against the entity. It returns a non-nil
// error if the token is invalid
func (c *UserEntity) Authenticate(session *mgo.Session, userAuth *UserAuthentication) error {
	// Check if credentails are not loaded then load the document.
	if !c.isCredentailsLoaded() {
		if err := common.MongoExecute(session, UserEntityDatabase, UserEntityCollection, func(co *mgo.Collection) error {
			return co.Find(bson.M{"public_id": userAuth.PublicID}).One(c)
		}); err != nil {
			return err
		}
	}

	return c.AuthenticateAgainst(userAuth.Token)
}

// AuthenticateAgainst authenticates the token against the entity. It returns a non-nil
// error if the token is invalid
func (c *UserEntity) AuthenticateAgainst(token string) error {
	if token == "" {
		return errors.New(credentailsAuthError)
	}

	return crypto.IsValidTokenForEntity(c, token)
}

// IsPasswordValid validates if the password belongs to the user entity and returns
// a non-nil error if the password is not a match or if the entity is in a
// invalid state.
func (c *UserEntity) IsPasswordValid(pwd string) error {
	if !c.isCredentailsLoaded() {
		return errors.New(invalidCrendentailsError)
	}

	hash, err := crypto.BcryptHash(c.PrivateID + pwd)
	if err != nil {
		return err
	}

	if !crypto.CompareBcryptHash(byte(c.Password), hash) {
		return errors.New("Invalid Password")
	}

	return nil
}

// SerializeAsPublic sets the company entity secret fields to defaults, to ensure
// only public allowed fields can be serialized.
func (c *UserEntity) SerializeAsPublic(includeMeta ...bool) {
	c.ID = ""
	c.Status = 0
	c.PrivateID = ""
	if len(includeMeta) == 0 {
		c.Config = nil
		c.ModifiedAt = nil
		c.CreatedAt = nil
	}
}

// SetToken sets the entity's token.
func (c *UserEntity) SetToken() error {
	token, err := crypto.TokenForEntity(c)
	if err != nil {
		return err
	}

	c.Token = base64.StdEncoding.EncodeToString(token)
	return nil
}

// Pwd returns the password to be used in creating the authentication
// token for the given entity. It satifies the crypto.Entity interface.
func (c *UserEntity) Pwd() ([]byte, error) {
	if c.PrivateID == "" {
		return nil, errors.New("Invalid Tenant Pwd")
	}

	return []byte(c.Password), nil
}

// Salt returns the salt key to be used in creating the authentication
// token for the given entity. It satifies the crypto.Entity interface.
func (c *UserEntity) Salt() ([]byte, error) {
	if c.PrivateID == "" || c.PublicID == "" {
		return nil, errors.New("Invalid Tenant Salt")
	}

	salt := fmt.Sprintf("%s:%s:%s", c.PublicID, c.PrivateID, fmt.Sprintf("%v", c.CreatedAt.UTC()))
	return salt, nil
}

// isCredentailsLoaded returns true/false stating the presence/absence of needed
// user entity fields.
func (c *UserEntity) isCredentailsLoaded() bool {
	if c.PublicID == "" || c.PrivateID == "" || c.Password == "" || c.Email == "" || c.Username == "" {
		return false
	}

	return true
}

// ValidatePassword validates if a password has a set of requirements.
func ValidatePassword(pwd string) error {
	if len(pw) < MinPasswordLength || len(pw) > MaxPasswordLength {
		return errors.New(InvalidPasswordError)
	}
	var num, lower, upper, spec bool
	for _, r := range pw {
		switch {
		case unicode.IsDigit(r):
			num = true
		case unicode.IsUpper(r):
			upper = true
		case unicode.IsLower(r):
			lower = true
		case unicode.IsSymbol(r), unicode.IsPunct(r):
			spec = true
		}
	}
	if num && lower && upper && spec {
		return nil
	}
	return errors.New(InvalidPasswordError)
}
