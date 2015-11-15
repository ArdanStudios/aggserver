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
	"github.com/ArdanStudios/aggserver/auth/vendor/github/satori/go.uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// CompanyEntityDatabase defines the database to use.
const CompanyEntityDatabase = ""

// CompanyEntityCollection defines the collection to use in storing company entites.
const CompanyEntityCollection = "companies"

// CompanyEntity represents a company based entity
// to provide authentication using a company token
type CompanyEntity struct {
	ID         bson.ObjectId          `json:"id" bson:"id"`
	Name       string                 `json:"name,omitempty" bson:"name"`
	Token      string                 `json:"token,omitempty" bson:"-"`
	PublicID   string                 `json:"public_id,omitempty" bson:"public_id"`
	PrivateID  string                 `json:"private_id,omitempty" bson:"private_id"`
	Status     EntityStatus           `json:"status" bson:"status"`
	Config     map[string]interface{} `json:"config,omitempty" bson:"config"`
	ModifiedAt *time.Time             `json:"modified_at,omitempty" bson:"modified_at"`
	CreatedAt  *time.Time             `json:"created_at,omitempty" bson:"created_at"`
}

// GetCompanyEntities returns a lists of all available company entities.
func GetCompanyEntities(session *mgo.Session) ([]*CompanyEntity, error) {
	var entities []*CompanyEntity

	if err := common.MongoExecute(session, CompanyEntityDatabase, CompanyEntityCollection, func(co *mgo.Collection) error {
		return co.Find(nil).All(&entity)
	}); err != nil {
		return nil, err
	}

	return entities, nil
}

// GetCompanyEntity returns a entity using the provided name.
func GetCompanyEntity(session *mgo.Session, name string) (*CompanyEntity, error) {
	var entity CompanyEntity

	if err := common.MongoExecute(session, CompanyEntityDatabase, CompanyEntityCollection, func(co *mgo.Collection) error {
		return co.Find(bson.M{"name": name}).One(&entity)
	}); err != nil {
		return nil, err
	}

	return &entity, nil
}

// NewCompanyEntity creates and initializes a ne CompanyEntity.
func NewCompanyEntity(name string, config map[string]interface{}) (*CompanyEntity, error) {
	if config == nil {
		config = make(map[string]interface{})
	}

	publicUUID, err := uuid.NewV4()
	privateUUID, err := uuid.NewV4()

	createdAt := time.Now()
	modifiedAt := time.Now()

	company := CompanyEntity{
		ID:         bson.NewObjectId(),
		Name:       name,
		Config:     config,
		PublicID:   publicUUID.String(),
		PrivateID:  privateUUID.String(),
		ModifiedAt: modifiedAt,
		CreatedAt:  createdAt,
	}

	// Setup the entity's authentication token
	if err := company.SetToken(); err != nil {
		return err
	}

	return &company
}

// Insert inserts the Company entity into the mongoDB database collection.
// It returns a non-nil error if an error occurs
func (c *CompanyEntity) Insert(session *mgo.Session) error {
	return common.MongoExecute(session, CompanyEntityDatabase, CompanyEntityCollection, func(co *mgo.Collection) error {
		return co.Insert(c)
	})
}

// CompanyNew provides a struct for use in creating a new company entity.
type CompanyNew struct {
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
}

// Create defines a new company entity and saves it into the giving entity database
// using the provided mongo session and serializable data.
func (c *CompanyEntity) Create(session *mgo.Session, data []byte) error {
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
		return nil, err
	}

	// Insert this record into the database.
	return common.MongoExecute(session, CompanyEntityDatabase, CompanyEntityCollection, func(co *mgo.Collection) error {
		return co.Insert(c)
	})
}

// CompanyUpdate provides a struct for use updating an existing company entity.
type CompanyUpdate struct {
	PublicID string                 `json:"public_id"`
	Name     string                 `json:"name"`
	Status   EntityStatus           `json:"status" `
	Config   map[string]interface{} `json:"config"`
}

// Update defines an update to a company's entity and updates the giving entity
// in the corresponding mongodb database and giving collection using the provided
// mongo session and serializable data.
func (c *CompanyEntity) Update(session *mgo.Session, data []byte) error {
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

	return common.MongoExecute(session, CompanyEntityDatabase, CompanyEntityCollection, func(co *mgo.Collection) error {
		return co.UpdateId(c.ID, c)
	})
}

// Destroy destroys/removes a company entity from the corresponding mongodb database
// and giving collection using the provided mongo session and serializable data.
func (c *CompanyEntity) Destroy(session *mgo.Session) error {
	return common.MongoExecute(session, CompanyEntityDatabase, CompanyEntityCollection, func(co *mgo.Collection) error {
		return co.RemoveId(c.ID)
	})
}

// Authenticate authenticates the token against the entity. It returns a non-nil
// error if the token is invalid
func (c *CompanyEntity) Authenticate(token string) error {
	if token == "" {
		return errors.New(credentailsAuthError)
	}

	return crypto.IsValidTokenForEntity(c, token)
}

// SerializeAsPublic sets the company entity secret fields to defaults, to ensure
// only public allowed fields can be serialized.
func (c *CompanyEntity) SerializeAsPublic(includeMeta ...bool) {
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
func (c *CompanyEntity) SetToken() error {
	token, err := crypto.TokenForEntity(c)
	if err != nil {
		return er
	}

	c.Token = base64.StdEncoding.EncodeToString(token)
	return nil
}

// Pwd returns the password to be used in creating the authentication
// token for the given entity. It satifies the crypto.Entity interface.
func (c *CompanyEntity) Pwd() ([]byte, error) {
	if c.PrivateID == "" {
		return nil, errors.New("Invalid Tenant Pwd")
	}

	return []byte(c.PrivateID), nil
}

// Salt returns the salt key to be used in creating the authentication
// token for the given entity. It satifies the crypto.Entity interface.
func (c *CompanyEntity) Salt() ([]byte, error) {
	if c.PrivateID == "" || c.PublicID == "" {
		return nil, errors.New("Invalid Tenant Salt")
	}

	salt := fmt.Sprintf("%s:%s:%s", c.PublicID, c.PrivateID, fmt.Sprintf("%v", c.CreatedAt.UTC()))
	return salt, nil
}

// ConfigKey loads the giving configuration key into the provided object else
// returns a non-nil error if the key does not exists or was not able to load
// into the giving store.
func (c *CompanyEntity) ConfigKey(key string, store interface{}) error {

	return nil
}
