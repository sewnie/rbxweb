// clientsettings selectively implements the 'clientsettings' Roblox Web API.
package clientsettings

import (
	"net/url"

	"github.com/apprehensions/rbxweb/internal/api"
)

// GamesServiceV1 handles the 'clientsettings/v2' Roblox Web API.
type ClientSettingsServiceV1 api.Service

// GetClientVersion gets the client version information for the named
// BinaryType and deployment channel.
func (c *ClientSettingsServiceV1) GetClientVersion(bt BinaryType, channel string) (*ClientVersion, error) {
	var cv ClientVersion

	path := c.Client.Path("v2/client-version/%s", nil, bt)
	if channel != "" {
		path += "/channel/" + channel
	}

	err := c.Client.Execute("GET", "clientsettings", path, nil, &cv)
	if err != nil {
		return nil, err
	}

	return &cv, nil
}

// GetUserChannel returns the channel name for the currently logged in
// user. The BinaryType given is optional; the Web API defaults to an unknown
// BinaryType.
func (c *ClientSettingsServiceV1) GetUserChannel(bt *BinaryType) (*UserChannel, error) {
	var uc UserChannel
	q := url.Values{}

	if bt != nil {
		q.Add("binaryType", string(*bt))
	}

	err := c.Client.Execute("GET",
		"clientsettings", c.Client.Path("/v2/user-channel", q), nil, &uc)
	if err != nil {
		return nil, err
	}

	return &uc, nil
}
