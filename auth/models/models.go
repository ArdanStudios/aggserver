// Package models contains data structures and associated behavior.
package models

import "github.com/coralproject/shelf/xenia/vendor/gopkg.in/bluesuncorp/validator.v6"

// Invalid describes a validation error belonging to a specific field.
type Invalid struct {
	Fld string `json:"field_name"`
	Err string `json:"error"`
}

var validate *validator.Validate

func init() {
	config := validator.Config{
		TagName:         "validate",
		ValidationFuncs: validator.BakedInValidators,
	}

	validate = validator.New(config)
}
