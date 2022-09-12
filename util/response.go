package util

import "fmt"

// An error raised when the data can not be found in the database
var ErrDataNotFound = fmt.Errorf("ErrDataNotFound")

var ErrInternalServer = fmt.Errorf("ErrInternalServer")

// A collection of authroization error messages that are used by the JWT auth middleware
var ErrAuthorization = fmt.Errorf("ErrAuthorization")

// An error message specific to indicate that the JWT has expired
var ErrTokenExpired = fmt.Errorf("ErrTokenExpired")

// A collection of validation error messages that are used by the validation middleware
var ErrValidation = fmt.Errorf("ErrValidation")

// An error message when the product path is not valid, path should be /products/[id]
var ErrInvalidRequestPath = fmt.Errorf("ErrInvalidRequestPath")

type GormError struct {
	Number  int    `json:"Number"`
	Message string `json:"Message"`
}

// SuccessResponse is returned when server successfully process the client request
type SuccessResponse struct {
	Message string `json:"message"`
}

// ErrorResponse is returned when errors occur
type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

// ErrorBody contains the detail of the error that will be returned inside the ErrorResponse object
type ErrorBody struct {
	// Server-defined set of error codes
	Code string `json:"code"`

	// A human-readable representation of the error
	Message string `json:"message"`

	// The target of the error
	Target string `json:"target"`
}
