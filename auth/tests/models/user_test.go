// Package models implements tests for the user API
package models

import (
	"testing"

	"github.com/ArdanStudios/aggserver/auth/models"
	"github.com/ArdanStudios/aggserver/auth/tests"
)

var user = new(models.User)

var u = models.UserNew{
	UserType:        1,
	FirstName:       "Josh",
	LastName:        "Zheng",
	Email:           "zheng@gmail.com",
	Company:         "Zuff",
	Password:        "Zhu-fro8bzr",
	PasswordConfirm: "Zhu-fro8bzr",
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

	userCreate(t)
	userUpdate(t)
	userLogin(t)
	userAuthenticate(t)
}

// userCreate tests the creation API for adding new users
func userCreate(t *testing.T) {
	t.Log("Given the need to create a new user.")
	{
		t.Log("\tWhen giving a models.UserNew struct")
		{
			err := user.Create(&u)
			if err != nil {
				t.Errorf("\t\tShould create user without errors %s", tests.Failed)
			} else {
				t.Logf("\t\tShould create user without errors %s", tests.Succeed)
			}

			if err := user.IsPasswordValid("Zhu-fro8bzr"); err != nil {
				t.Errorf("\t\tShould have a valid password %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have a valid password %s", tests.Succeed)
			}

			if _, err := user.Validate(); err != nil {
				t.Errorf("\t\tShould be a valid user without errors %s", tests.Failed)
			} else {
				t.Logf("\t\tShould be a valid user without errors %s", tests.Succeed)
			}

		}
	}
}

// userUpdating tests the user API for the removal of users
func userUpdate(t *testing.T) {
	t.Log("Given the need to update a user.")
	{
		t.Log("\tWhen giving a models.UserUpdate struct")
		{
			err := user.Update(&models.UserUpdate{
				FirstName: "William",
				PublicID:  user.PublicID,
			})

			if err != nil {
				t.Errorf("\t\tShould have updated user without errors %s", tests.Failed)
			} else {
				t.Logf("\t\tShould have updated user without errors %s", tests.Succeed)
			}

			if user.FirstName != "William" {
				t.Errorf("\t\tShould have first name as %q %s", "William", tests.Failed)
			} else {
				t.Logf("\t\tShould have first name as %q %s", "William", tests.Succeed)
			}
		}
	}
}

// userLogin tests the authentication process users using the password and
// email crendentials
func userLogin(t *testing.T) {
	t.Log("Given the need to login with a user.")
	{
		t.Log("\tWhen giving user email and password credentials")
		{
			err := user.AuthenticateLogin(&models.UserLoginAuthentication{
				Email:    "zheng@gmail.com",
				Password: "Zhu-fro8bzr",
			})

			if err != nil {
				t.Errorf("\t\tShould successfully authenticate %s", tests.Failed)
			} else {
				t.Logf("\t\tShould successfully authenticate %s", tests.Succeed)
			}
		}
	}
}

// userAuthenticate tests the authentication process users using the token and
// publicID crendentials
func userAuthenticate(t *testing.T) {
	t.Log("Given the need to authenticate with a user.")
	{
		t.Log("\tWhen giving user token and public_id credentials")
		{

			err := user.AuthenticateToken(&models.UserTokenAuthentication{
				PublicID: user.PublicID,
				Token:    user.Token,
			})

			if err != nil {
				t.Errorf("\t\tShould successfully authenticate %s", tests.Failed)
			} else {
				t.Logf("\t\tShould successfully authenticate %s", tests.Succeed)
			}

		}
	}
}
