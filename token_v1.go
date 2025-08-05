package rbxweb

// AuthTokenServiceV1 partially handles the undocumented 'auth-token-service/v1' Roblox Web API.
type AuthTokenServiceV1 service

// Token is a representation an unknown model returned by login/create.
type Token struct {
	Code           string `json:"code"`
	Status         string `json:"status"`
	PrivateKey     string `json:"privateKey"`
	ExpirationTime string `json:"expirationTime"`
	ImagePath      string `json:"imagePath"`
}

// TokenStatus is a representation an unknown model returned by login/status.
type TokenStatus struct {
	Status            string `json:"status"`
	AccountName       string `json:"accountName"`
	AccountPictureURL string `json:"accountPictureUrl"`
	ExpirationTime    string `json:"expirationTime"`
}

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

	if err := a.Client.csrfRequired(); err != nil {
		return nil, err
	}

	err := a.Client.Execute("POST", "apis", "auth-token-service/v1/login/status", req, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
