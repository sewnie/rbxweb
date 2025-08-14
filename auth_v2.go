package rbxweb

// AuthServiceV2 partially handles the 'auth/v2' Roblox Web API.
type AuthServiceV2 service

// LoginType represents an available credentials type for LoginRequest
type LoginType string

const (
	LoginTypeUsername = "Username"
	LoginTypeToken    = "AuthToken"
)

// Login implements the LoginRequest API model.
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

// CreateLogin logins as the user with the given Token.
//
// If logging in with a username and password, set value and password to a
// username and password, and force the request to be downgraded to HTTP/1.1
// by providing the Client with a custom transport. This is unconfirmed.
//
// If logging in with a Token, set the login type to LoginTypeToken, the value to
// the token's code, and the password to the token's private key.
func (a *AuthServiceV2) CreateLogin(value, password string, login LoginType) (*Login, error) {
	lreq := struct {
		CType    string `json:"ctype"`
		CValue   string `json:"cvalue"`
		Password string `json:"password"`
	}{
		CType:    string(login),
		CValue:   value,
		Password: password,
	}

	req, err := a.Client.NewRequest("POST", "auth", "v2/login", lreq)
	if err != nil {
		return nil, err
	}

	// This can definitely be revoked by Roblox if they care to do so.
	// Once that happens, it will mean rbxweb will have to initialize a tracker
	// when required for a request.
	//
	// req.AddCookie is not used in this endpoint since net/http sanitizes and
	// changes the final cookie value.
	//
	// Unsure if it is required.
	req.Header.Set("Cookie", "RBXEventTrackerV2=CreateDate=08/01/2025 12:38:07&rbxid=227740497&browserid=1748902424900004")

	l := new(Login)
	if _, err := a.Client.Do(req, &l); err != nil {
		return nil, err
	}

	return l, nil
}
