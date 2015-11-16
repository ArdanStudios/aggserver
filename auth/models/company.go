package models

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ArdanStudios/aggserver/auth/common"
	"github.com/ArdanStudios/aggserver/auth/crypto"
	"github.com/ArdanStudios/aggserver/auth/vendor/github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// CompanyDatabase defines the database to use.
const CompanyDatabase = ""

// CompanyCollection defines the collection to use in storing company entites.
const CompanyCollection = "companies"

// Company represents a company based entity
// to provide authentication using a company token
type Company struct {
	ID         bson.ObjectId          `json:"id" bson:"id"`
	Name       string                 `json:"name,omitempty" bson:"name"`
	Token      string                 `json:"token,omitempty" bson:"-"`
	PublicID   string                 `json:"public_id,omitempty" bson:"public_id"`
	PrivateID  string                 `json:"private_id,omitempty" bson:"private_id"`
	Status     common.EntityStatus    `json:"status" bson:"status"`
	Config     map[string]interface{} `json:"config,omitempty" bson:"config"`
	ModifiedAt *time.Time             `json:"modified_at,omitempty" bson:"modified_at"`
	CreatedAt  *time.Time             `json:"created_at,omitempty" bson:"created_at"`
}

// NewCompany creates and initializes a new Company.
func NewCompany(name string, config map[string]interface{}) (*Company, error) {
	if config == nil {
		config = make(map[string]interface{})
	}

	publicUUID, err := uuid.NewV4()
	privateUUID, err := uuid.NewV4()

	createdAt := time.Now()
	modifiedAt := time.Now()

	company := Company{
		ID:         bson.NewObjectId(),
		Name:       name,
		Config:     config,
		PublicID:   publicUUID.String(),
		PrivateID:  privateUUID.String(),
		ModifiedAt: modifiedAt,
		CreatedAt:  createdAt,
		Status:     common.ActiveStatus,
	}

	// Setup the entity's authentication token
	if err := company.SetToken(); err != nil {
		return err
	}

	return &company
}

// CompanyNew provides a struct for use in creating a new company entity.
type CompanyNew struct {
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
}

// Create defines a new company entity and saves it into the giving entity database
// using the provided mongo session and serializable data.
func (c *Company) Create(session *mgo.Session, data []byte) error {
	var newCompany CompanyNew

	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&newCompany); err != nil {
		return err
	}

	publicUUID := uuid.NewV4()
	privateUUID := uuid.NewV4()

	createdAt := time.Now()
	modifiedAt := time.Now()

	c.Name = newCompany.Name
	c.PublicID = publicUUID.String()
	c.PrivateID = privateUUID.String()
	c.CreatedAt = createdAt
	c.ModifiedAt = modifiedAt
	c.Status = EntityActive

	if newCompany.Config != nil {
		c.Config = newCompany.Config
	} else {
		c.Config = make(map[string]interface{})
	}

	// Set up the entity's authentication token.
	if err := c.SetToken(); err != nil {
		return err
	}

	return nil
}

// CompanyUpdate provides a struct for use updating an existing company entity.
type CompanyUpdate struct {
	PublicID string                 `json:"public_id"`
	Name     string                 `json:"name"`
	Status   common.EntityStatus    `json:"status" `
	Config   map[string]interface{} `json:"config"`
}

// Update defines an update to a company's entity and updates the giving entity
// in the corresponding mongodb database and giving collection using the provided
// mongo session and serializable data.
func (c *Company) Update(session *mgo.Session, data []byte) error {
	var updatingCompany CompanyUpdate

	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&updatingCompany); err != nil {
		return err
	}

	if updatingCompany.PublicID != c.PublicID {
		return errors.New("Invalid PublicID for entity")
	}

	c.Name = updatingCompany.Name
	c.Status = updatingCompany.Status

	if updatingCompany.Config != nil {
		c.Config = updatingCompany.Config
	}

	c.ModifiedAt = time.Now()

	return common.MongoExecute(session, CompanyDatabase, CompanyCollection, func(co *mgo.Collection) error {
		return co.UpdateId(c.ID, c)
	})
}

// Authenticate authenticates the token against the entity. It returns a non-nil
// error if the token is invalid
func (c *Company) Authenticate(token string) error {
	if token == "" {
		return errors.New(common.CredentailsAuthError)
	}

	return crypto.IsValidTokenForEntity(c, token)
}

// SerializeAsPublic sets the company entity secret fields to defaults, to ensure
// only public allowed fields can be serialized.
func (c *Company) SerializeAsPublic(includeMeta ...bool) {
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
func (c *Company) SetToken() error {
	token, err := crypto.TokenForEntity(c)
	if err != nil {
		return er
	}

	c.Token = base64.StdEncoding.EncodeToString(token)
	return nil
}

// Pwd returns the password to be used in creating the authentication
// token for the given entity. It satifies the crypto.Entity interface.
func (c *Company) Pwd() ([]byte, error) {
	if c.PrivateID == "" {
		return nil, errors.New("Invalid Tenant Pwd")
	}

	return []byte(c.PrivateID), nil
}

// Salt returns the salt key to be used in creating the authentication
// token for the given entity. It satifies the crypto.Entity interface.
func (c *Company) Salt() ([]byte, error) {
	if c.PrivateID == "" || c.PublicID == "" {
		return nil, errors.New("Invalid Tenant Salt")
	}

	salt := fmt.Sprintf("%s:%s:%s", c.PublicID, c.PrivateID, fmt.Sprintf("%v", c.CreatedAt.UTC()))
	return salt, nil
}

// ConfigKey loads the giving configuration key into the provided object else
// returns a non-nil error if the key does not exists or was not able to load
// into the giving store.
func (c *Company) ConfigKey(key string, store interface{}) error {

	return nil
}
