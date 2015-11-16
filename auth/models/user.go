package models

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"
	"unicode"

	"github.com/ArdanStudios/aggserver/auth/common"
	"github.com/ArdanStudios/aggserver/auth/crypto"
	"github.com/satori/go.uuid"

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

// UserAddress contains information about a user's address.
type UserAddress struct {
	Type         int        `bson:"type" json:"type" validate:"required"`
	LineOne      string     `bson:"line_one" json:"line_one" validate:"required"`
	LineTwo      string     `bson:"line_two" json:"line_two,omitempty"`
	City         string     `bson:"city" json:"city" validate:"required"`
	State        string     `bson:"state" json:"state" validate:"required"`
	Zipcode      string     `bson:"zipcode" json:"zipcode" validate:"required"`
	Phone        string     `bson:"phone" json:"phone" validate:"required"`
	DateModified *time.Time `bson:"date_modified" json:"date_modified"`
	DateCreated  *time.Time `bson:"date_created,omitempty" json:"date_created"`
}

// Validate checks the fields to verify the value is in a proper state.
func (ua *UserAddress) Validate() ([]Invalid, error) {
	var inv []Invalid

	errs := validate.Struct(ua)
	if errs != nil {
		for _, err := range errs {
			inv = append(inv, Invalid{Fld: err.Field, Err: err.Tag})
		}

		return inv, errors.New("Validation failures identified")
	}

	return nil, nil
}

// Compare checks the fields against another UserAddress value.
func (ua *UserAddress) Compare(uat *UserAddress) ([]Invalid, error) {
	var inv []Invalid

	if ua.Type != uat.Type {
		inv = append(inv, Invalid{Fld: "Type", Err: fmt.Sprintf("The value of Type is not the same. %d != %d", ua.Type, uat.Type)})
	}

	if ua.LineOne != uat.LineOne {
		inv = append(inv, Invalid{Fld: "LineOne", Err: fmt.Sprintf("The value of LineOne is not the same. %s != %s", ua.LineOne, uat.LineOne)})
	}

	if ua.City != uat.City {
		inv = append(inv, Invalid{Fld: "City", Err: fmt.Sprintf("The value of City is not the same. %s != %s", ua.City, uat.City)})
	}

	if ua.State != uat.State {
		inv = append(inv, Invalid{Fld: "State", Err: fmt.Sprintf("The value of State is not the same. %s != %s", ua.State, uat.State)})
	}

	if ua.Zipcode != uat.Zipcode {
		inv = append(inv, Invalid{Fld: "Zipcode", Err: fmt.Sprintf("The value of Zipcode is not the same. %s != %s", ua.Zipcode, uat.Zipcode)})
	}

	if ua.Phone != uat.Phone {
		inv = append(inv, Invalid{Fld: "Phone", Err: fmt.Sprintf("The value of Phone is not the same. %s != %s", ua.Phone, uat.Phone)})
	}

	if len(inv) > 0 {
		return inv, errors.New("Compare failures identified")
	}

	return nil, nil
}

// User represents an entity for user records.
type User struct {
	ID         bson.ObjectId       `json:"id" bson:"id"`
	UserType   common.UserType     `json:"user_type" bson:"user_type" validate:"required"`
	FirstName  string              `json:"first_name,omitempty" bson:"first_name" validate:"required"`
	LastName   string              `json:"last_name,omitempty" bson:"last_name" validate:"required"`
	Company    string              `json:"company,omitempty" bson:"company" validate:"required"`
	Addresses  []UserAddress       `bson:"addresses" json:"addresses" validate:"required"`
	Token      string              `json:"token,omitempty" bson:"-"`
	PublicID   string              `json:"public_id,omitempty" bson:"public_id" validate:"required"`
	PrivateID  string              `json:"private_id,omitempty" bson:"private_id"`
	Email      string              `bson:"email" json:"email" validate:"required"`
	Status     common.EntityStatus `json:"status" bson:"status" validate:"required"`
	Password   string              `bson:"password" json:"password,omitempty" validate:"required"`
	ModifiedAt *time.Time          `json:"modified_at,omitempty" bson:"modified_at"`
	CreatedAt  *time.Time          `json:"created_at,omitempty" bson:"created_at"`
}

// Validate checks the fields to verify the value is in a proper state.
func (u *User) Validate() ([]Invalid, error) {
	var inv []Invalid

	errs := validate.Struct(u)
	if errs != nil {
		for _, err := range errs {
			inv = append(inv, Invalid{Fld: err.Field, Err: err.Tag})
		}

		return inv, errors.New("Validation failures identified")
	}

	for _, ua := range u.Addresses {
		if va, err := ua.Validate(); err != nil {
			inv = append(inv, va...)
		}
	}

	if len(inv) > 0 {
		return inv, errors.New("Validation failures identified")
	}

	return nil, nil
}

// Compare checks the fields against another User value.
func (u *User) Compare(ut *User) ([]Invalid, error) {
	var inv []Invalid

	if u.UserType != ut.UserType {
		inv = append(inv, Invalid{Fld: "Type", Err: fmt.Sprintf("The value of Type is not the same. %d != %d", u.UserType, ut.UserType)})
	}

	if u.FirstName != ut.FirstName {
		inv = append(inv, Invalid{Fld: "FirstName", Err: fmt.Sprintf("The value of FirstName is not the same. %s != %s", u.FirstName, ut.FirstName)})
	}

	if u.LastName != ut.LastName {
		inv = append(inv, Invalid{Fld: "LastName", Err: fmt.Sprintf("The value of LastName is not the same. %s != %s", u.LastName, ut.LastName)})
	}

	if u.Email != ut.Email {
		inv = append(inv, Invalid{Fld: "Email", Err: fmt.Sprintf("The value of Email is not the same. %s != %s", u.Email, ut.Email)})
	}

	if u.Company != ut.Company {
		inv = append(inv, Invalid{Fld: "Company", Err: fmt.Sprintf("The value of Company is not the same. %s != %s", u.Company, ut.Company)})
	}

	uLen := len(u.Addresses)
	utLen := len(ut.Addresses)

	if uLen != utLen {
		inv = append(inv, Invalid{Fld: "Addresses", Err: fmt.Sprintf("The set of Addresses is not the same. %d != %d", uLen, utLen)})
	}

	for idx, ua := range u.Addresses {
		if idx >= utLen {
			break
		}

		if va, err := ua.Compare(&ut.Addresses[idx]); err != nil {
			inv = append(inv, va...)
		}
	}

	if len(inv) > 0 {
		return inv, errors.New("Compare failures identified")
	}

	return nil, nil
}

// UserNew provides a struct for use in creating a new User entity.
type UserNew struct {
	UserType        common.UserType `json:"user_type"`
	FirstName       string          `json:"first_name"`
	LastName        string          `json:"last_name"`
	Email           string          `json:"email"`
	Company         string          `json:"company"`
	Addresses       []UserAddress   `json:"addresses"`
	Password        string          `json:"password"`
	PasswordConfirm string          `json:"password_confirm"`
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

// Create defines a new User entity and saves it into the giving entity database
// using the provided mongo session and serializable data.
func (u *User) Create(newUser *UserNew) error {
	if err := newUser.ValidatePassword(); err != nil {
		return err
	}

	publicUUID := uuid.NewV4()
	privateUUID := uuid.NewV4()

	createdAt := time.Now()
	modifiedAt := time.Now()

	u.FirstName = newUser.FirstName
	u.LastName = newUser.LastName
	u.UserType = newUser.UserType
	u.Company = newUser.Company
	u.Email = newUser.Email
	u.PublicID = publicUUID.String()
	u.PrivateID = privateUUID.String()
	u.CreatedAt = &createdAt
	u.ModifiedAt = &modifiedAt
	u.Status = common.ActiveEntity

	// Add the address.
	u.Addresses = append(u.Addresses, newUser.Addresses...)

	// Create the user password hash.
	p, err := crypto.BcryptHash([]byte(u.PrivateID + newUser.Password))
	if err != nil {
		return err
	}

	u.Password = p

	// Set up the entity's authentication token.
	if err := u.SetToken(); err != nil {
		return err
	}

	return nil
}

// UserPasswordReset is used to resets a users entity's password.
type UserPasswordReset struct {
	ID       bson.ObjectId `bson:"id" json:"id,omitempty"`
	PublicID string        `bson:"public_id" json:"public_id"`
	Token    string        `bson:"token" json:"token"`
	ExpireAt *time.Time    `bson:"expire_at" json:"expire_at"`
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

// ChangePassword changes the password according to the values being received.
func (u *User) ChangePassword(changeUser *UserPasswordChange) error {
	if !u.isCredentailsLoaded() {
		return errors.New(common.InvalidCredentailsError)
	}

	p, err := crypto.BcryptHash([]byte(u.PrivateID + changeUser.Password))
	if err != nil {
		return err
	}

	u.Password = string(p)

	ms := time.Now()
	u.ModifiedAt = &ms

	if err := u.SetToken(); err != nil {
		return err
	}

	return nil
}

// UserUpdate provides a struct for use updating an existing User entity.
type UserUpdate struct {
	FirstName string              `json:"first_name"`
	LastName  string              `json:"last_name"`
	Addresses []UserAddress       `bson:"addresses" json:"addresses"`
	PublicID  string              `json:"public_id"`
	Email     string              `json:"email"`
	Status    common.EntityStatus `json:"status" `
	Token     string              `json:"token"`
	UserType  common.UserType     `json:"user_type"`
	Company   string              `json:"company"`
}

// Update defines an update to a User's entity and updates the giving entity
// in the corresponding mongodb database and giving collection using the provided
// mongo session and serializable data.
func (u *User) Update(updatingUser *UserUpdate) error {
	// Check if credentails are not loaded then load the document.
	if !u.isSafeCredentailsLoaded() {
		return errors.New(common.InvalidCredentailsError)
	}

	if updatingUser.PublicID != u.PublicID {
		return errors.New("Invalid PublicID for entity")
	}

	if updatingUser.FirstName != "" {
		u.FirstName = updatingUser.FirstName
	}

	if updatingUser.LastName != "" {
		u.LastName = updatingUser.LastName
	}

	if len(updatingUser.Addresses) > 0 {
		u.Addresses = append([]UserAddress{}, updatingUser.Addresses...)
	}

	if updatingUser.Email != "" {
		u.Email = updatingUser.Email
	}

	if updatingUser.Company != "" {
		u.Company = updatingUser.Company
	}

	if updatingUser.Status != common.NoStatusEntity {

		u.Status = updatingUser.Status
	}

	if updatingUser.UserType != common.UnknownType {
		u.UserType = updatingUser.UserType
	}

	ms := time.Now()
	u.ModifiedAt = &ms

	return nil
}

// UserDestroy contains the given user public_id to instruct a removal of the entity
// from the databse.
type UserDestroy struct {
	PublicID string `json:"public_id"`
}

// UserLoginAuthentication provides the necessary information needed for authenticating
// a user entity using the login crendentails.
type UserLoginAuthentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthenticateLogin authenticates the login crendentials against the given user
// entity. It returns a non-nil error if the crendetails are invalid.
func (u *User) AuthenticateLogin(userLogin *UserLoginAuthentication) error {
	// Check if credentails are not loaded then load the document.
	if !u.isCredentailsLoaded() {
		return errors.New(common.InvalidCredentailsError)
	}

	loginHash := []byte(u.PrivateID + userLogin.Password)
	// Generate the hash corresponding with the given password
	passwordHash, err := crypto.BcryptHash(loginHash)
	if err != nil {
		return err
	}

	return u.AuthenticateAgainst(passwordHash)
}

// UserTokenAuthentication provides the necessary information needed for authenticating
// a user entity using the Token and PublicID credentials.
type UserTokenAuthentication struct {
	PublicID string `json:"public_id"`
	Token    string `json:"public_id"`
}

// AuthenticateToken authenticates the token against the entity. It returns a non-nil
// error if the token is invalid
func (u *User) AuthenticateToken(userAuth *UserTokenAuthentication) error {
	// Check if credentails are not loaded then load the document.
	if !u.isCredentailsLoaded() {
		return errors.New(common.InvalidCredentailsError)
	}

	if u.PublicID != userAuth.PublicID {
		return errors.New(common.CredentailsAuthError)
	}

	return u.AuthenticateAgainst(userAuth.Token)
}

// AuthenticateAgainst authenticates the token against the entity. It returns a non-nil
// error if the token is invalid
func (u *User) AuthenticateAgainst(token string) error {
	if token == "" {
		return errors.New(common.CredentailsAuthError)
	}

	return crypto.IsTokenValidForEntity(u, token)
}

// IsPasswordValid validates if the password belongs to the user entity and returns
// a non-nil error if the password is not a match or if the entity is in a
// invalid state.
func (u *User) IsPasswordValid(pwd string) error {
	if !u.isCredentailsLoaded() {
		return errors.New(common.InvalidCredentailsError)
	}

	hash, err := crypto.BcryptHash([]byte(u.PrivateID + pwd))
	if err != nil {
		return err
	}

	binHash, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return err
	}

	userBinHash, err := base64.StdEncoding.DecodeString(u.Password)
	if err != nil {
		return err
	}

	if err := crypto.CompareBcryptHash(binHash, userBinHash); err != nil {
		return err
	}

	return nil
}

// SerializeAsPublic sets the User entity secret fields to defaults, to ensure
// only public allowed fields can be serialized.
func (u *User) SerializeAsPublic(includeMeta ...bool) {
	u.ID = ""
	u.Status = 0
	u.PrivateID = ""
	if len(includeMeta) == 0 {
		u.ModifiedAt = nil
		u.CreatedAt = nil
	}
}

// SetToken sets the entity's token.
func (u *User) SetToken() error {
	token, err := crypto.TokenForEntity(u)
	if err != nil {
		return err
	}

	u.Token = base64.StdEncoding.EncodeToString(token)
	return nil
}

// Pwd returns the password to be used in creating the authentication
// token for the given entity. It satifies the crypto.Entity interface.
func (u *User) Pwd() ([]byte, error) {
	if u.PrivateID == "" {
		return nil, errors.New("Invalid Tenant Pwd")
	}

	return []byte(u.Password), nil
}

// Salt returns the salt key to be used in creating the authentication
// token for the given entity. It satifies the crypto.Entity interface.
func (u *User) Salt() ([]byte, error) {
	if u.PrivateID == "" || u.PublicID == "" {
		return nil, errors.New("Invalid Tenant Salt")
	}

	salt := fmt.Sprintf("%s:%s:%s", u.PublicID, u.PrivateID, fmt.Sprintf("%v", u.CreatedAt.UTC()))
	return []byte(salt), nil
}

// isSafeCredentailsLoaded returns true/false stating the presence/absence of the
// basic user entity fields.
func (u *User) isSafeCredentailsLoaded() bool {
	if u.PublicID == "" || u.Password == "" || u.Email == "" || u.UserType == 0 {
		return false
	}

	return true
}

// isCredentailsLoaded returns true/false stating the presence/absence of needed
// user entity fields.
func (u *User) isCredentailsLoaded() bool {
	if u.PublicID == "" || u.PrivateID == "" || u.Password == "" || u.Email == "" {
		return false
	}

	return true
}

// ValidatePassword validates if a password has a set of requirements.
func ValidatePassword(pw string) error {
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
