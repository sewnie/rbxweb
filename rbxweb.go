// Package rbxweb provides API routines to interact with Roblox's web API.
package rbxweb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Client embeds an [http.Client], used to make Roblox API requests.
//
// BaseDomain is the URL domain used to execute calls to, in case an alternative
// domain is given.
type Client struct {
	http.Client
	BaseDomain string

	Security string // .ROBLOSECURITY
	Token    string // X-CSRF-Token

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	GamesV1          *GamesServiceV1
	ThumbnailsV1     *ThumbnailsServiceV1
	UsersV1          *UsersServiceV1
	AuthV2           *AuthServiceV2
	OAuthV1          *OAuthServiceV1
	ClientSettingsV2 *ClientSettingsServiceV2
	AuthTokenV1      *AuthTokenServiceV1
}

// NewClient returns a new Client.
func NewClient() *Client {
	c := &Client{
		BaseDomain: "roblox.com",
	}

	c.common.Client = c
	c.GamesV1 = (*GamesServiceV1)(&c.common)
	c.ThumbnailsV1 = (*ThumbnailsServiceV1)(&c.common)
	c.UsersV1 = (*UsersServiceV1)(&c.common)
	c.AuthV2 = (*AuthServiceV2)(&c.common)
	c.OAuthV1 = (*OAuthServiceV1)(&c.common)
	c.ClientSettingsV2 = (*ClientSettingsServiceV2)(&c.common)
	c.AuthTokenV1 = (*AuthTokenServiceV1)(&c.common)

	return c
}

type service struct {
	Client *Client
}

// path constructs a URL path with the given path as the format, values (if any),
// and format parameters for the path. The encoded query will be appended to the format.
func path(format string, query url.Values, a ...any) string {
	if query != nil {
		format += "?" + query.Encode()
	}
	return fmt.Sprintf(format, a...)
}

// NewRequest returns a new API request with the given relative path and
// the service (subdomain) to use with the BaseDomain of the Client. If a body
// is specified, and it is of type [url.Values], it will be added to the request
// as application/x-www-form-urlencoded, otherwise, the body is used as
// application/json if non-nil.
//
// The request returned expects a application/json.
//
// The security cookie and CSRF token will be added to the request if available.
func (c *Client) NewRequest(method, service, path string, body any) (*http.Request, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.BaseDomain,
		Path:   path,
	}
	if service != "" {
		u.Host = service + "." + u.Host
	}

	buf := new(bytes.Buffer)
	content := ""
	if v, ok := body.(url.Values); ok {
		buf.WriteString(v.Encode())
		content = "application/x-www-form-urlencoded"
	} else if body != nil {
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(body); err != nil {
			return nil, err
		}
		content = "application/json"
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "rbxweb/v0.0.0")

	if content != "" {
		req.Header.Set("Content-Type", content)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "identity")

	if c.Token != "" {
		req.Header.Set("X-CSRF-TOKEN", c.Token)
	}

	if c.Security != "" {
		req.AddCookie(&http.Cookie{
			Name:  ".ROBLOSECURITY",
			Value: c.Security,
		})
	}

	return req, nil
}

// Do performs the API request and returns the HTTP response. If any error occurs,
// the respose body will be closed If a API error response is available, it will be
// returned as either an Errors or string error for undocumented APIs; if all
// else fails, a StatusError will be returned. Otherwise, the user is responsible for
// handling and closing the response body.
//
// If the response returned a security cookie it will be used in future requests.
//
// If the request fails with 403 and returns X-CSRF-TOKEN, GetBody will be used from the
// request, as the underlying type made from [NewRequest] is bytes.Buffer, and the
// request will be tried again with the new X-CSRF-TOKEN, It will also be stored
// and used for future requests until the cycle occurs again.
func (c *Client) BareDo(req *http.Request) (*http.Response, error) {
	resp, err := c.Client.Do(req)
	if err != nil {
		return resp, err
	}

	t := resp.Header.Get("X-CSRF-TOKEN")
	if t != "" && resp.StatusCode == http.StatusForbidden {
		resp.Body.Close()
		c.Token = t

		req = req.Clone(req.Context())
		req.Header.Set("X-CSRF-TOKEN", c.Token)
		if req.GetBody != nil {
			req.Body, err = req.GetBody()
			if err != nil {
				return nil, err
			}
		}
		resp, err = c.Client.Do(req)
		if err != nil {
			return resp, err
		}
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == ".ROBLOSECURITY" {
			c.Security = cookie.Value
		}
	}

	// Skip reading for an error if the response is acceptable
	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return resp, nil
	}
	defer resp.Body.Close()

	content := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(content, "application/json") {
		return resp, &StatusError{StatusCode: resp.StatusCode}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, &StatusError{StatusCode: resp.StatusCode}
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	errResp := new(Errors)
	if err := dec.Decode(errResp); err == nil {
		return resp, fmt.Errorf("api errors: %w", errResp)
	}

	// Some undocumented APIs return a single string as an error
	var errStr string
	if err := json.Unmarshal(data, &errStr); err == nil {
		return resp, errors.New(errStr)
	}

	return resp, fmt.Errorf("unhandled error of %d: %s",
		resp.StatusCode, string(data))
}

// Do performs the API request and returns the HTTP response and decodes
// or writes the response to v, if non-nil, as necessary.
// The response body of the HTTP request is always going to be closed.
// See [BareDo] for more details.
func (c *Client) Do(req *http.Request, v any) (*http.Response, error) {
	resp, err := c.BareDo(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, err
}

// Executes wraps around NewRequest and Do for immediate execution of a request.
func (c *Client) Execute(method, service, path string, body any, v any) error {
	req, err := c.NewRequest(method, service, path, body)
	if err != nil {
		return err
	}

	if _, err := c.Do(req, v); err != nil {
		return err
	}

	return nil
}

// StatusError represents an unexpected HTTP error, in the case
// that a ErrorResponse was unable to be parsed.
type StatusError struct {
	StatusCode int
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("bad response: %s", http.StatusText(e.StatusCode))
}

// Error implements the error response model of the API.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

// errorsResponse implements the errors response model of the API.
type Errors struct {
	Errors []Error `json:"errors,omitempty"`
}

// Error implements the error interface.
func (err Error) Error() string {
	return fmt.Sprintf("response code %d: %s", err.Code, err.Message)
}

// Error implemements the error interface.
func (errs Errors) Error() string {
	s := make([]string, len(errs.Errors))
	for i, e := range errs.Errors {
		s[i] = e.Error()
	}
	return strings.Join(s, "; ")
}

// Unwrap implements the Unwrap interface by returning the first error in the
// list.
func (errs Errors) Unwrap() error {
	if len(errs.Errors) == 0 {
		return nil
	}
	return errs.Errors[0]
}

func formatSlice[T any](values []T) []string {
	if len(values) == 0 {
		return nil
	}

	s := make([]string, len(values))
	for i, v := range values {
		s[i] = fmt.Sprintf("%v", v)
	}
	return s
}

func getList[T any](v []T, err error) (*T, error) {
	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, nil
	}
	return &v[0], nil
}
