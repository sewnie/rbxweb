package rbxweb

import (
	"fmt"
	"strings"
)

// ErrorResponse implements the error response model of the API.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

// errorsResponse implements the errors response model of the API.
type errorsResponse struct {
	Errors []ErrorResponse `json:"errors,omitempty"`
}

// Error implements the error interface.
func (err ErrorResponse) Error() string {
	return fmt.Sprintf("response code %d: %s", err.Code, err.Message)
}

// Error implemements the error interface.
func (errs errorsResponse) Error() string {
	s := make([]string, len(errs.Errors))
	for i, e := range errs.Errors {
		s[i] = e.Error()
	}
	return strings.Join(s, "; ")
}

// Unwrap implements the Unwrap interface by returning the first error in the
// list.
func (errs errorsResponse) Unwrap() error {
	if len(errs.Errors) == 0 {
		return nil
	}
	return errs.Errors[0]
}
