package util

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator"
)

// ValidationError wraps the validators FieldError so we do not
// expose this to out code
type ValidationError struct {
	validator.FieldError
}

func (v ValidationError) Error() string {
	return fmt.Sprintf(
		"Key: '%s' Error: Field validation for '%s' failed on the '%s' tag",
		v.Namespace(),
		v.Field(),
		v.Tag(),
	)
}

// ValidationErrors is a collection of ValidationError
type ValidationErrors []ValidationError

// Errors converts the slice into a string slice
func (v ValidationErrors) Errors() []string {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

type Validation struct {
	validate *validator.Validate
}

// NewValidation creates a new Validation type
func NewValidation() *Validation {
	validate := validator.New()
	// validate.RegisterValidation("sku", validateSKU)

	return &Validation{validate}
}

func (v *Validation) Validate(i interface{}) ValidationErrors {
	errs := v.validate.Struct(i)

	if errs != nil {
		// Asserts that errs is not nil and that the value stored in errs is of type ValidationErrors
		validationErrs := errs.(validator.ValidationErrors)

		var returnErrs []ValidationError
		for _, err := range validationErrs {
			// cast the FieldError into our ValidationError and append to the slice
			ve := ValidationError{err}
			returnErrs = append(returnErrs, ve)
		}

		return returnErrs
	}
	return nil
}

func validateSKU(fl validator.FieldLevel) bool {
	// SKU must be in format abc-absd-dfsdf
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	return len(matches) == 1
}
