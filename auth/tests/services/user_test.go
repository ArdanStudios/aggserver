package services

import (
	"log"
	"testing"

	"github.com/ArdanStudios/aggserver/auth/session"
	"github.com/coralproject/shelf/xenia/tests"
)

var newUser = `{
	"user_type":  1,
	"first_name": "John",
	"last_name":  "Row",
	"email":      "john.row@gmail.com",
	"company":    "Rowling Bank",
	"address": [{
			"type":     1,
			"line_one": "12973 Lane ST",
			"line_two": "FUMI 153",
			"city":     "Zhigi",
			"state":    "FL",
			"zipcode":  "53172",
			"phone":    "0808-629-4323",
	}]
}`

// TestUsers tests the models.User API.
func TestUsers(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	dbSession, err := session.Session()
	if err != nil {
		log.Fatal("Unable to retrieve database session")
	}

	t.Log("Given the need to create a new user.")
	{
		t.Log("When giving a models.UserNew struct")
		{

			t.Log("Should create user without errors")
		}
	}

}
