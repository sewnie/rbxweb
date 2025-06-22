package rbxweb

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrBadStatus = errors.New("bad status")
	ErrNoData    = errors.New("no data")
)

// StatusError represents an unexpected HTTP error, in the case
// that a ErrorResponse was unable to be parsed.
type StatusError struct {
	StatusCode int
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("bad response: %s", http.StatusText(e.StatusCode))
}

// ErrorResponse implements the error response model of the API.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

// errorsResponse implements the errors response model of the API.
type ErrorsResponse struct {
	Errors []ErrorResponse `json:"errors,omitempty"`
}

// Error implements the error interface.
func (err ErrorResponse) Error() string {
	return fmt.Sprintf("response code %d: %s", err.Code, err.Message)
}

// Error implemements the error interface.
func (errs ErrorsResponse) Error() string {
	s := make([]string, len(errs.Errors))
	for i, e := range errs.Errors {
		s[i] = e.Error()
	}
	return strings.Join(s, "; ")
}

// Unwrap implements the Unwrap interface by returning the first error in the
// list.
func (errs ErrorsResponse) Unwrap() error {
	if len(errs.Errors) == 0 {
		return nil
	}
	return errs.Errors[0]
}
