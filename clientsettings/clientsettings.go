// clientsettings implements the Roblox ClientSettings Web API.
package clientsettings

import (
	"net/url"

	"github.com/apprehensions/rbxweb"
)

// BinaryType represents a supported Roblox platform.
type BinaryType string

const (
	WindowsPlayer   BinaryType = "WindowsPlayer"
	WindowsStudio   BinaryType = "WindowsStudio" // Deprecated in favor of WindowsStudio64, undocumented
	WindowsStudio64 BinaryType = "WindowsStudio64"
	MacPlayer       BinaryType = "MacPlayer"
	MacStudio       BinaryType = "MacStudio"
)

// Short returns the shortened form of BinaryType.
func (bt BinaryType) Short() string {
	switch bt {
	case WindowsPlayer, MacPlayer:
		return "Player"
	case WindowsStudio, WindowsStudio64, MacStudio:
		return "Studio"
	default:
		return "ShortenedBinaryType(" + string(bt) + ")"
	}
}

// ClientVersion implements the ClientVersionResponse API model.
type ClientVersion struct {
	Version      string `json:"version"`
	GUID         string `json:"clientVersionUpload"`
	Bootstrapper string `json:"bootstrapperVersion"`
	NextGUID     string `json:"nextClientVersionUpload,omitempty"`
	NextVersion  string `json:"nextClientVersion,omitempty"`
}

// String implements the Stringer interface.
func (cv ClientVersion) String() string {
	return cv.Version + " (" + cv.GUID + ")"
}

// AssignmentType represents how the user was bound to a channel.
type AssignmentType string

const (
	None                  AssignmentType = "None"
	PerMille              AssignmentType = "PerMille"
	BoundToPrivateChannel AssignmentType = "BoundToPrivateChannel"
	BoundToPublicChannel  AssignmentType = "BoundToPublicChannel"
)

// UserChannel implements the UserChannelResponse API model.
type UserChannel struct {
	Channel    string         `json:"channelName"`
	Assignment AssignmentType `json:"channelAssignmentType"`
	Token      string         `json:"token"`
}

// GetClientVersion gets the client version information for the named
// BinaryType and deployment channel, using the 'clientsettingscdn' as oppose
// to the 'clientsetting' API provider.
func GetClientVersion(bt BinaryType, channel string) (*ClientVersion, error) {
	var cv ClientVersion

	ep := "v2/client-version/" + string(bt)
	if channel != "" {
		ep += "/channel/" + channel
	}

	err := rbxweb.Request("GET", rbxweb.GetURL("clientsettings", ep, nil), nil, &cv)
	if err != nil {
		return nil, err
	}

	return &cv, nil
}

// GetUserChannel returns the channel name for the currently logged in
// user. The BinaryType given is optional; the Web API defaults to an unknown
// BinaryType.
func GetUserChannel(bt *BinaryType) (*UserChannel, error) {
	var uc UserChannel
	q := url.Values{}

	if bt != nil {
		q.Add("binaryType", string(*bt))
	}

	err := rbxweb.Request("GET",
		rbxweb.GetURL("clientsettings", "/v2/user-channel", q), nil, &uc)
	if err != nil {
		return nil, err
	}

	return &uc, nil
}
