package auth

// LoginType represents an available credentials type for LoginRequest
type LoginType string

const (
	Username = "Username"
	Token    = "AuthToken"
)

type Login struct {
	User struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
	} `json:"user"`
	TwoStepVerificationData struct {
		MediaType int    `json:"mediaType"`
		Ticket    string `json:"ticket"`
	} `json:"twoStepVerificationData"`
	IdentityVerificationLoginTicket string `json:"identityVerificationLoginTicket"`
	IsBanned                        bool   `json:"isBanned"`
	AccountBlob                     string `json:"accountBlob"`
	ShouldUpdateEmail               bool   `json:"shouldUpdateEmail"`
	RecoveryEmail                   string `json:"recoveryEmail"`
}

type loginIntent struct {
	PublicKey    string `json:"clientPublicKey"`
	Epoch        int64  `json:"clientEpochTimestamp"`
	ServerNonce  string `json:"serverNonce"`
	SaiSignature string `json:"saiSignature"`
}

type loginRequest struct {
	CType    string      `json:"ctype"`
	CValue   string      `json:"cvalue"`
	Password string      `json:"password"`
	Intent   loginIntent `json:"secureAuthenticationIntent"`
}
