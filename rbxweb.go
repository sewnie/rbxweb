// Package rbxweb provides API routines to interact with Roblox's web API.
package rbxweb

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

const (
	cookieSecurity  = ".ROBLOSECURITY"
	headerCSRFToken = "x-csrf-token"
)

// Client embeds an [http.Client], used to make Roblox API requests.
//
// BaseDomain is the URL domain used to execute calls to, in case an alternative
// domain is given.
type Client struct {
	http.Client
	BaseDomain string

	// Only for debugging purposes.
	Logger *slog.Logger

	Security  *http.Cookie
	csrfToken string

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	GamesV1          *GamesServiceV1
	ThumbnailsV1     *ThumbnailsServiceV1
	UsersV1          *UsersServiceV1
	AuthV2           *AuthServiceV2
	ClientSettingsV2 *ClientSettingsServiceV2
	AuthTokenV1      *AuthTokenServiceV1
}

// NewClient returns a new Client.
func NewClient() *Client {
	c := &Client{
		BaseDomain: "roblox.com",

		Client: http.Client{Transport: &http.Transport{
			// Fixes authentication endpoints
			TLSClientConfig: &tls.Config{},
		}},
	}

	c.common.Client = c
	c.GamesV1 = (*GamesServiceV1)(&c.common)
	c.ThumbnailsV1 = (*ThumbnailsServiceV1)(&c.common)
	c.UsersV1 = (*UsersServiceV1)(&c.common)
	c.AuthV2 = (*AuthServiceV2)(&c.common)
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
// is specified, it will be added to the request body as JSON.
// The security cookie and CSRF token will be added to the request if available.
func (c *Client) NewRequest(method, service, path string, body any) (*http.Request, error) {
	url := url.URL{
		Scheme: "https",
		Host:   c.BaseDomain,
		Path:   path,
	}
	if service != "" {
		url.Host = service + "." + url.Host
	}

	buf := new(bytes.Buffer)
	if body != nil {
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(body); err != nil {
			return nil, err
		}
		buf.Truncate(buf.Len() - 1)
	}

	c.logDebug("New Request",
		"service", service, "method", method, "body", buf.String())

	req, err := http.NewRequest(method, url.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "rbxweb/v0.0.0")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "identity")

	// > ... Similarly, RoundTrip should not attempt to
	// > handle higher-level protocol details such as redirects,
	// > authentication, or cookies.

	if c.csrfToken != "" {
		req.Header.Set(headerCSRFToken, c.csrfToken)
	}

	if c.Security != nil {
		req.AddCookie(c.Security)
	}

	return req, nil
}

// Do performs the API request and returns the HTTP response. If any error occurs,
// the respose body will be closed If a API error response is available, it will be
// returned as either an Errors or string error for undocumented APIs; if all
// else fails, a StatusError will be returned. Otherwise, the user is responsible for
// handling and closing the response body.
//
// If the response returned a security cookie or a X-CSRF-TOKEN header, it will
// be used in future requests. If a request rate limits or returns a header for
// resending the request, it will be returned as is.
func (c *Client) BareDo(req *http.Request) (*http.Response, error) {
	resp, err := c.Client.Do(req)
	if err != nil {
		return resp, err
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == cookieSecurity {
			c.logDebug("Recieved " + cookieSecurity)
			c.Security = cookie
		}
	}

	if t := resp.Header.Get(headerCSRFToken); t != "" {
		c.csrfToken = t
		c.logDebug("Recieved CSRF", "token", c.csrfToken)
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

	errResp := new(Errors)
	if err := json.Unmarshal(data, errResp); err == nil {
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
		c.logDebug("Response", "status", resp.StatusCode,
			"content", resp.Header.Get("Content-Type"))
		_, err = io.Copy(v, resp.Body)
	default:
		err = json.NewDecoder(resp.Body).Decode(v)
		c.logDebug("Response", "status", resp.StatusCode, "data", v)
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

// rbxweb does not automatically retry a request if it requires a XSRF token, instead
// endpoints that require this must use it beforehand for easier API usage.
// in the future, automatically using the recieved XSRF token upon a "XSRF token invalid"
// may be used if necessary.
func (c *Client) csrfRequired() error {
	if c.csrfToken != "" {
		return nil
	}
	return c.AuthV2.setCSRFToken()
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
