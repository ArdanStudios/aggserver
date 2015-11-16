package services

import (
	"github.com/ArdanStudios/aggserver/auth/common"
	"github.com/ArdanStudios/aggserver/auth/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type companyService struct{}

// CompanyService provides a global handler for Company entity CRUD and authentication management.
var CompanyService companyService

// GetCompanyEntities returns a lists of all available company entities.
func (c *companyService) GetCompanies(session *mgo.Session) ([]*Company, error) {
	var entities []*models.Company

	if err := common.MongoExecute(session, CompanyDatabase, CompanyCollection, func(co *mgo.Collection) error {
		return co.Find(nil).All(&entities)
	}); err != nil {
		return nil, err
	}

	return entities, nil
}

// GetCompany returns a entity using the provided name.
func (c *companyService) GetCompanyByName(session *mgo.Session, name string) (*Company, error) {
	var entity models.Company

	if err := common.MongoExecute(session, CompanyDatabase, CompanyCollection, func(co *mgo.Collection) error {
		return co.Find(bson.M{"name": name}).One(&entity)
	}); err != nil {
		return nil, err
	}

	return &entity, nil
}
