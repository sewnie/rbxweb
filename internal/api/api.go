package api

import (
	"fmt"
	"net/http"
	"net/url"
)

type Client interface {
	BareDo(req *http.Request) (*http.Response, error)
	Do(req *http.Request, v any) (*http.Response, error)
	NewRequest(method, service, path string, body any) (*http.Request, error)
	Execute(method, service, path string, body any, v any) error
}

type Service struct {
	Client Client
}

func GetList[T any](v []T, err error) (*T, error) {
	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, nil
	}
	return &v[0], nil
}

// Path constructs a URL path with the given path as the format, values (if any),
// and format parameters for the path. The encoded query will be appended to the format.
func Path(format string, query url.Values, a ...any) string {
	if query != nil {
		format += "?" + query.Encode()
	}
	return fmt.Sprintf(format, a...)
}
