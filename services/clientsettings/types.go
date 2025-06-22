package clientsettings

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
