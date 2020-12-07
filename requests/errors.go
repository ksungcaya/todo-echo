package requests

import (
	"encoding/json"
)

// ValidationErrors struct
type ValidationErrors struct {
	Errors map[string][]string `json:"errors"`
}

// ResponseError struct
type ResponseError struct {
	Errors interface{} `json:"errors"`
}

// NewValidationError creates a ValidationError.
func NewValidationError(field string, message string) ValidationErrors {
	e := make(map[string][]string, 1)
	e[field] = []string{message}
	return ValidationErrors{Errors: e}
}

// NewResponseError creates a ResponseError with "errors" as parent node.
func NewResponseError(err error) ResponseError {
	e := ResponseError{}
	switch v := err.(type) {
	case ValidationErrors:
		e.Errors = v.Errors
	default:
		e.Errors = map[string]string{"message": v.Error()}
	}
	return e
}

// It allows ValidationErrors to subscribe to the Error interface.
// The error map can be accessed through ve.Errors field.
func (ve ValidationErrors) Error() string {
	data, _ := json.MarshalIndent(ve.Errors, "", "  ")
	return string(data)
}
