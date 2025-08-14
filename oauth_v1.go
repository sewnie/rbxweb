package rbxweb

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strconv"
)

// OauthServiceV1 partially handles the 'oauth/v1' Roblox Web API.
type OAuthServiceV1 service

// OAuthClientID represents a OAuth client ID.
type OAuthClientID string // for convenience

// PermissionScope implements an unknown API model for OAuth permission management.
type PermissionScope struct {
	// One of "openid", "profile", "email", "verification", "credentials", "age", "premium", "roles".
	Type string `json:"scopeType"`
	// One of "read", "write".
	Operations []string `json:"operations"`
}

// PermissionResourceOwner implements an unknown API model for a OAuth resource owner.
type PermissionResourceOwner struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// PermissionResourceInfo implements an unknown API model for OAuth resource information.
type PermissionResourceInfo struct {
	Owner     PermissionResourceOwner `json:"owner"`
	Resources struct{}                `json:"resources"` // Unknown, required
}

type OAuthToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`      // space-separated scopes
	TokenType    string `json:"token_type"` // "Bearer"
}

// GetToken uses undocumented parts of oauth/v1/token to get OAuth authentication
// for Roblox Studio
func (o *OAuthServiceV1) AuthStudioToken(c OAuthClientID, u *AuthStudioURL) (*OAuthToken, error) {
	q := url.Values{}
	q.Add("code", u.URL.Query().Get(("code")))
	q.Add("grant_type", "authorization_code")
	q.Add("client_id", string(c))
	q.Add("code_verifier", u.Verifier)

	t := new(OAuthToken)
	err := o.Client.Execute("POST", "apis", "oauth/v1/token", q, &t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func NewCodeVerifierBytes(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

type AuthStudioURL struct {
	URL      *url.URL
	Verifier string
}

// GetAuthorizations uses the undocumented oauth/v1/authorizations endpoint to return
// a roblox-studio-auth scheme URL to be used for authenticating Roblox Studio.
// The returned URL should be used with [AuthStudioToken].
func (o *OAuthServiceV1) GetAuthStudioURL(c OAuthClientID, userID UserID) (*AuthStudioURL, error) {
	codeRaw := make([]byte, 32)
	_, err := rand.Read(codeRaw)
	if err != nil {
		return nil, err
	}
	code := base64.RawURLEncoding.EncodeToString(codeRaw)

	h := sha256.Sum256([]byte(code))
	challenge := base64.RawURLEncoding.EncodeToString(h[:])

	stateRaw := map[string]string{
		"random_string": code,
		"pid":           "220",
	}
	stateJSON, _ := json.Marshal(stateRaw)
	state := base64.RawURLEncoding.EncodeToString(stateJSON)

	data := struct {
		ClientID      string                   `json:"clientId"` // Retrieved from Studio OAuth2Config.json
		Challenge     string                   `json:"codeChallenge"`
		Method        string                   `json:"codeChallengeMethod"`
		Nonce         string                   `json:"nonce"`
		ResponseTypes []string                 `json:"responseTypes"` // One of "Code", "None"
		RedirectURI   string                   `json:"redirectUri"`
		Scopes        []PermissionScope        `json:"scopes"`
		State         string                   `json:"state"`
		ResourceInfos []PermissionResourceInfo `json:"resourceInfos"`
	}{
		ClientID:      string(c),
		Challenge:     challenge,
		Method:        "S256",
		Nonce:         "id-roblox",
		ResponseTypes: []string{"Code"},
		RedirectURI:   "roblox-studio-auth:/",
		Scopes: []PermissionScope{
			{Type: "openid", Operations: []string{"read"}},
			{Type: "credentials", Operations: []string{"read"}},
			{Type: "profile", Operations: []string{"read"}},
			{Type: "age", Operations: []string{"read"}},
			{Type: "roles", Operations: []string{"read"}},
			{Type: "premium", Operations: []string{"read"}},
		},
		State: state,
		ResourceInfos: []PermissionResourceInfo{
			{
				Owner: PermissionResourceOwner{
					ID:   strconv.FormatInt(int64(userID), 10),
					Type: "User",
				},
				Resources: struct{}{},
			},
		},
	}

	respData := struct {
		Location string `json:"location"`
	}{}

	err = o.Client.Execute("POST", "apis", "oauth/v1/authorizations", data, &respData)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(respData.Location)
	if err != nil {
		return nil, err
	}

	return &AuthStudioURL{
		URL:      u,
		Verifier: code,
	}, nil
}
