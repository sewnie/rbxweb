// token partially implements the undocumented auth-token-service API.
package token

import (
	"github.com/apprehensions/rbxweb/internal/api"
)

// ThumbnailsServiceV1 partially handles the undocumented 'auth-token-service/v1' Roblox Web API.
type AuthTokenServiceV1 api.Service

// CreateToken returns a newly created token.
func (a *AuthTokenServiceV1) CreateToken() (*Token, error) {
	var t Token

	err := a.Client.Execute("POST", "apis", "auth-token-service/v1/login/create", nil, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// GetTokenStatus returns the status of a Token.
func (a *AuthTokenServiceV1) GetTokenStatus(t *Token) (*TokenStatus, error) {
	var s TokenStatus
	req := struct {
		Code       string `json:"code"`
		PrivateKey string `json:"privateKey"`
	}{t.Code, t.PrivateKey}

	err := a.Client.Execute("POST", "apis", "auth-token-service/v1/login/status", req, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
