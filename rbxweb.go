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

// Path constructs a URL path with the given path as the format, values (if any),
// and format parameters for the path. The encoded query will be appended to the format.
func Path(format string, query url.Values, a ...any) string {
	if query != nil {
		format += "?" + query.Encode()
	}
	return fmt.Sprintf(format, a...)
}

// BareDo will execute the given HTTP request, leaving the response body
// to be read by the user. If any error occurs, the respose body will be closed.
// If a API error response is available, it will be returned as ErrorsResponse
// or string for undocumented APIs, otherwise, a StatusError will be returned.
//
// If the response returned a security cookie or a X-CSRF-TOKEN header, it will
// be used in future requests. If a request that failed returns this header, the
// request will not be re-attempted.
func (c *Client) BareDo(req *http.Request) (*http.Response, error) {
	c.logInfo("Performing Request",
		"method", req.Method, "path", req.URL.Path)

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
// encoded and will be added to the request body. The security cookie will be added
// to the request if available.
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

func (c *Client) logIsDebug() bool {
	return c.Logger != nil && c.Logger.Enabled(nil, slog.LevelDebug)
}
