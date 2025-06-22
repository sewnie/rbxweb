package api

import (
	"net/http"
	"net/url"
)

type Client interface {
	BareDo(req *http.Request) (*http.Response, error)
	Do(req *http.Request, v any) (*http.Response, error)
	NewRequest(method, service, path string, body any) (*http.Request, error)
	Execute(method, service, path string, body any, v any) error
	Path(format string, query url.Values, a ...any) string
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
