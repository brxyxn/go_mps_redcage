package data

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"testing"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"
)

func TestValidValues(t *testing.T) {
	p := Client{
		Firstname: "John",
		Lastname:  "Doe",
		Username:  "jdoe",
	}

	re := regexp.MustCompile(`([A-Za-z])+`)
	extractedFirstName := re.FindString(p.Firstname)
	extractedLastName := re.FindString(p.Lastname)

	assert.Equal(t, extractedFirstName, "John")
	assert.NotEqual(t, extractedFirstName, "john")
	assert.NotEmpty(t, extractedFirstName, "Shouldn't be empty")

	assert.Equal(t, extractedLastName, "Doe")
	assert.NotEqual(t, extractedLastName, "doe")
	assert.NotEmpty(t, extractedLastName, "Shouldn't be empty")

	re = regexp.MustCompile(`([a-z]{3,})`)
	extractedUserName := re.FindString(p.Username)

	assert.Equal(t, extractedUserName, "jdoe")
	assert.NotEqual(t, extractedUserName, "do")
	assert.NotEmpty(t, extractedUserName, "Shouldn't be empty")
}

func TestClientToJSON(t *testing.T) {
	ps := []*Client{
		{
			Firstname: "John",
			Lastname:  "Doe",
			Username:  "jdoe",
		},
	}

	b := bytes.NewBufferString("")
	err := ToJSON(ps, b)
	assert.NoError(t, err)
}

/*

Validation stuff

*/

func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

type ValidationError struct {
	validator.FieldError
}

type ValidationErrors []ValidationError

func (v *Validation) Validate(i interface{}) ValidationErrors {
	errs := v.validate.Struct(i).(validator.ValidationErrors)

	if len(errs) == 0 {
		return nil
	}

	var returnErrs []ValidationError
	for _, err := range errs {
		// cast the FieldError into our ValidationError and append to the slice
		ve := ValidationError{err.(validator.FieldError)}
		returnErrs = append(returnErrs, ve)
	}

	return returnErrs
}

type Validation struct {
	validate *validator.Validate
}

func NewValidation() *Validation {
	validate := validator.New()
	validate.RegisterValidation("firstname", validateFirstName)

	return &Validation{validate}
}

// validate firstname
func validateFirstName(fl validator.FieldLevel) bool {
	// firstname format: John
	re := regexp.MustCompile(`([A-Za-z])+`)
	name_len := re.FindString(fl.Field().String())

	if len(name_len) == 1 {
		return true
	}

	return false
}
