// Package models implements tests for the user API
package models

import (
	"log"
	"testing"

	"github.com/ArdanStudios/aggserver/auth/models"
	"github.com/ArdanStudios/aggserver/auth/session"
	"github.com/coralproject/shelf/xenia/tests"
	"gopkg.in/mgo.v2"
)

var user = new(models.User)

var u = models.UserNew{
	Type:      1,
	FirstName: "Josh",
	LastName:  "Zheng",
	Email:     "zheng@gmail.com",
	Company:   "Zuff",
	Addresses: []models.UserAddress{
		{
			Type:    1,
			LineOne: "12973 Lane ST",
			LineTwo: "FUMI 153",
			City:    "Zhigi",
			State:   "FL",
			Zipcode: "53172",
			Phone:   "0808-629-4323",
		},
	},
}

// TestUsers tests the models.User API.
func TestUsers(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	dbSession, err := session.Session()
	if err != nil {
		log.Fatal("Unable to retrieve database session")
	}

	userCreate(t, dbSession)
	userInsert(t, dbSession)
	userUpdate(t, dbSession)
	userDestroy(t, dbSession)
	userLogin(t, dbSession)
	userAuthenticate(t, dbSession)
}

// userCreate tests the creation API for adding new users
func userCreate(t *testing.T, session *mgo.Session) {
	t.Log("Given the need to create a new user.")
	{
		t.Log("When giving a models.UserNew struct")
		{

			t.Log("Should create user without errors")
		}
	}
}

// userInsert tests the user API for adding new users into the database
func userInsert(t *testing.T, session *mgo.Session) {
	t.Log("Given the need to create a new user.")
	{
		t.Log("When giving a models.User struct")
		{

		}
	}
}

// userUpdating tests the user API for the removal of users
func userUpdate(t *testing.T, session *mgo.Session) {
	t.Log("Given the need to create a new user.")
	{
		t.Log("When giving a models.UserUpdate struct")
		{

		}
	}
}

// userDestroy tests the user API for the removal of users
func userDestroy(t *testing.T, session *mgo.Session) {
	t.Log("Given the need to create a new user.")
	{
		t.Log("When giving a user public_id")
		{

		}
	}
}

// userLogin tests the authentication process users using the password and
// email crendentials
func userLogin(t *testing.T, session *mgo.Session) {
	t.Log("Given the need to create a new user.")
	{
		t.Log("When giving user email and password credentials")
		{

		}
	}
}

// userAuthenticate tests the authentication process users using the token and
// publicID crendentials
func userAuthenticate(t *testing.T, session *mgo.Session) {
	t.Log("Given the need to create a new user.")
	{
		t.Log("When giving user token and public_id credentials")
		{

		}
	}
}
