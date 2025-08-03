// Package rbxweb provides API routines to interact with Roblox's web API.
//
//go:generate go run ./cmd/genservices
package rbxweb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/sewnie/rbxweb/internal/api"
)

const (
	cookieSecurity = ".ROBLOSECURITY"
	headerToken    = "x-csrf-token"
)

// Client embeds an [http.Client], used to make Roblox API requests.
//
// BaseDomain is the URL domain used to execute calls to, in case an alternative
// domain is given.
type Client struct {
	http.Client
	BaseDomain string

	Logger *slog.Logger

	Security string
	Token    string

	common api.Service // Reuse a single struct instead of allocating one for each service on the heap.

	Services
}

// NewClient returns a new Client.
func NewClient() *Client {
	c := &Client{
		BaseDomain: "roblox.com",
	}

	c.common.Client = c
	c.setServices()

	return c
}

// BareDo will execute the given HTTP request, leaving the response body
// to be read by the user. If any error occurs, the respose body will be closed.
// If a API error response is available, it will be returned as ErrorsResponse
// or string for undocumented APIs, otherwise, a StatusError will be returned.
//
// If Roblox set a ROBLOSECURITY cookie or a X-CSRF-TOKEN header, it will
// always be used in future requests.
func (c *Client) BareDo(req *http.Request) (*http.Response, error) {
	c.logInfo("Performing Request",
		"method", req.Method, "path", req.URL.Path)

	resp, err := c.Client.Do(req)
	if err != nil {
		return resp, err
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == cookieSecurity {
			c.Security = cookie.Value
		}
	}
	if t := resp.Header.Get(headerToken); t != "" {
		c.Token = t
	}
	if resp.StatusCode == http.StatusOK {
		return resp, nil
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, &StatusError{StatusCode: resp.StatusCode}
	}

	errResp := new(ErrorsResponse)
	if err := json.Unmarshal(data, errResp); err == nil {
		return resp, fmt.Errorf("api errors: %w", errResp)
	}

	// Some undocumented APIs return a single string as an error
	var errStr string
	if err := json.Unmarshal(data, &errStr); err == nil {
		return resp, errors.New(errStr)
	}

	return resp, &StatusError{StatusCode: resp.StatusCode}
}

// Do will perform the given HTTP request on the client, and will write
// the response body as necessary to v, if non-nil. In any case, the response
// will always be returned as it is with a resulting error, if any.
// The response body of the HTTP request is always going to be closed.
func (c *Client) Do(req *http.Request, v any) (*http.Response, error) {
	resp, err := c.BareDo(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	switch v := v.(type) {
	case nil:
		return resp, nil
	case io.Writer:
		if c.logIsDebug() {
			b, _ := io.ReadAll(resp.Body)
			c.logDebug("Response", "status", resp.StatusCode, "data", string(b))
			v.Write(b)
			return resp, nil
		}
		_, err = io.Copy(v, resp.Body)
		return resp, err
	default:
		err = json.NewDecoder(resp.Body).Decode(v)
		c.logDebug("Response", "status", resp.StatusCode, "data", v)
		return resp, err
	}
}

// NewRequest returns a new http.Request, with a path that will be relative
// to the BaseDomain of the client, and a service - which can be empty, to indicate
// the microservice to use. If body is specified, it will be interepreted as JSON
// encoded and will be added to the request body.
func (c *Client) NewRequest(method, service, path string, body any) (*http.Request, error) {
	url := url.URL{
		Scheme: "https",
		Host:   c.BaseDomain,
		Path:   path,
	}
	if service != "" {
		url.Host = service + "." + url.Host
	}

	c.logDebug("New Request",
		"service", service, "method", method, "body", body)

	buf := new(bytes.Buffer)
	if body != nil {
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url.String(), buf)
	if err != nil {
		return nil, err
	}

	if c.Security != "" {
		req.AddCookie(&http.Cookie{
			Name:  cookieSecurity,
			Value: c.Security,
		})
	}
	if c.Token != "" {
		req.Header.Set(headerToken, c.Token)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// Executes creates a new Request with NewRequest and the given parameters and
// immediately executes it with Do, which unmarshals the response body to v, if any.
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

func (c *Client) logDebug(msg string, args ...any) {
	if c.Logger != nil {
		c.Logger.Debug("rbxweb: "+msg, args...)
	}
}

func (c *Client) logInfo(msg string, args ...any) {
	if c.Logger != nil {
		c.Logger.Info("rbxweb: "+msg, args...)
	}
}

func (c *Client) logWarn(msg string, args ...any) {
	if c.Logger != nil {
		c.Logger.Warn("rbxweb: "+msg, args...)
	}
}

func (c *Client) logError(msg string, args ...any) {
	if c.Logger != nil {
		c.Logger.Error("rbxweb: "+msg, args...)
	}
}

func (c *Client) logIsDebug() bool {
	return c.Logger != nil && c.Logger.Enabled(nil, slog.LevelDebug)
}
