package rbxweb

import (
	"net/url"
)

// ClientSettingsServiceV2 partially handles the 'clientsettings/v2' Roblox Web API.
type ClientSettingsServiceV2 service

// BinaryType represents a supported Roblox platform.
type BinaryType string

const (
	BinaryTypeWindowsPlayer   BinaryType = "WindowsPlayer"
	BinaryTypeWindowsStudio   BinaryType = "WindowsStudio" // Deprecated in favor of WindowsStudio64, undocumented
	BinaryTypeWindowsStudio64 BinaryType = "WindowsStudio64"
	BinaryTypeMacPlayer       BinaryType = "MacPlayer"
	BinaryTypeMacStudio       BinaryType = "MacStudio"
)

// Short returns the shortened form of BinaryType.
func (bt BinaryType) Short() string {
	switch bt {
	case BinaryTypeWindowsPlayer, BinaryTypeMacPlayer:
		return "Player"
	case BinaryTypeWindowsStudio, BinaryTypeWindowsStudio64, BinaryTypeMacStudio:
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
	AssignmentTypeNone     AssignmentType = "None"
	AssignmentTypePerMille AssignmentType = "PerMille"
	AssignmentTypePrivate  AssignmentType = "BoundToPrivateChannel"
	AssignmentTypePublic   AssignmentType = "BoundToPublicChannel"
)

// UserChannel implements the UserChannelResponse API model.
type UserChannel struct {
	Channel    string         `json:"channelName"`
	Assignment AssignmentType `json:"channelAssignmentType"`
	Token      string         `json:"token"`
}

// GetClientVersion gets the client version information for the named
// BinaryType and deployment channel.
func (c *ClientSettingsServiceV2) GetClientVersion(bt BinaryType, channel string) (*ClientVersion, error) {
	var cv ClientVersion

	path := Path("v2/client-version/%s", nil, bt)
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
func (c *ClientSettingsServiceV2) GetUserChannel(bt *BinaryType) (*UserChannel, error) {
	var uc UserChannel
	q := url.Values{}

	if bt != nil {
		q.Add("binaryType", string(*bt))
	}

	err := c.Client.Execute("GET",
		"clientsettings", Path("/v2/user-channel", q), nil, &uc)
	if err != nil {
		return nil, err
	}

	return &uc, nil
}
