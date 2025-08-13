package rbxweb

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strconv"
)

// OauthServiceV1 partially handles the 'oauth/v1' Roblox Web API.
type OAuthServiceV1 service

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

// GetAuthorizations implements the undocumented oauth/v1/authorizations endpoint,
// it is only being used to return a roblox-studio-auth URI until more documentation
// or uses are found.
func (o *OAuthServiceV1) GetAuthorizations(clientID string, userID UserID) (string, error) {
	codeRaw := make([]byte, 32)
	_, err := rand.Read(codeRaw)
	if err != nil {
		return "", err
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
		ClientID:      clientID,
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

	if err := o.Client.csrfRequired(); err != nil {
		return "", err
	}

	err = o.Client.Execute("POST", "apis", "oauth/v1/authorizations", data, &respData)
	if err != nil {
		return "", err
	}

	return respData.Location, nil
}
