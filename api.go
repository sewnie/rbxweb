// Package rbxweb provides API routines to interact with Roblox's web API.
package rbxweb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// Client is the [http.Client] used to make Roblox API requests.
var Client = &http.Client{}

var (
	ErrBadStatus = errors.New("bad status")
	ErrNoData    = errors.New("no data")
)

// Request performs a Roblox API request given a method, url returned by
// [GetURL], body and data interfaces to use as data to send and to recieve.
func Request(method, url string, body, data interface{}) error {
	buf := new(bytes.Buffer)
	if body != nil {
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		// Return API error only if such error exists
		e := new(errorsResponse)
		if err := dec.Decode(e); err == nil {
			return e
		}
		return fmt.Errorf("%w: %s", ErrBadStatus, resp.Status)
	}

	if data != nil {
		return dec.Decode(data)
	}

	return nil
}

// GetURL constructs a Roblox web API URL with the given service as the
// subdomain, and the given arguments, with HTTPS as the protocol.
func GetURL(service string, path string, query url.Values) string {
	url := url.URL{
		Scheme: "https",
		Host:   "roblox.com",
		Path:   path,
	}
	if query != nil {
		url.RawQuery = query.Encode()
	}
	if service != "" {
		url.Host = service + "." + url.Host
	}
	return url.String()
}
