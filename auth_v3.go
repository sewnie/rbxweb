package rbxweb

// AuthServiceV3 partially handles the 'auth/v3' Roblox Web API.
type AuthServiceV3 service

type TwoStepLogin struct {
	Ticket      string `json:"identityVerificationLoginTicket"`
	AccountBlob string `json:"accountBlob"`
}

// CreateLoginTwoStep logins as the user with the given TwoStepToken.
func (a *AuthServiceV3) CreateTwoStepLogin(rememberDevice bool, token *TwoStepToken) (*TwoStepLogin, error) {
	reqData := struct {
		*TwoStepToken
		RememberDevice bool `json:"rememberDevice"`
	}{
		TwoStepToken:   token,
		RememberDevice: rememberDevice,
	}
	l := new(TwoStepLogin)

	if err := a.Client.Execute("POST", "auth",
		path("v3/users/%d/two-step-verification/login", nil, token.User), reqData, l); err != nil {
		return nil, err
	}

	return l, nil
}
