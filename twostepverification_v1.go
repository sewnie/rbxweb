package rbxweb

// TwoStepVerificationServiceV1 partially handles the 'twostepverification/v1' Roblox Web API.
type TwoStepVerificationServiceV1 service

// VerificationType should represent the TwoStepVerificationSentResponse API model,
// But the API internally uses a string and not an int.
type VerificationType string

const (
	VerificationTypeEmail         = "Email"
	VerificationTypeSMS           = "SMS"
	VerificationTypeAuthenticator = "Authenticator"
	VerificationTypeRecoveryCode  = "RecoveryCode"
	VerificationTypeSecurityKey   = "SecurityKey"
	VerificationTypeCrossDevice   = "CrossDevice"
	VerificationTypePassword      = "Password"
	VerificationTypePasskey       = "Passkey"
)

// v1/users/{userid}/challenges/{map.Key}/verify
var verificationPaths = map[VerificationType]string{
	VerificationTypeEmail:         "email",
	VerificationTypeSMS:           "sms",
	VerificationTypeAuthenticator: "authenticator",
	VerificationTypeRecoveryCode:  "recovery-codes",
	VerificationTypeSecurityKey:   "security-key",
	VerificationTypeCrossDevice:   "cross-device",
	VerificationTypePassword:      "password",
	VerificationTypePasskey:       "passkey",
}

// TwoStepToken is used to store two step verification token data.
type TwoStepToken struct {
	User      UserID
	Token     string `json:"verificationToken"`
	Challenge string `json:"challengeId"`
}

// CreateToken creates a two step verification token for the current Login and given code.
//
// The type of code to send is determined by [Login.Verification.Type].
func (t *TwoStepVerificationServiceV1) CreateToken(l *Login, code string) (*TwoStepToken, error) {
	reqData := struct {
		ActionType string `json:"actionType"`
		Challenge  string `json:"challengeId"`
		Code       string `json:"code"`
	}{
		ActionType: "Login",
		Challenge:  l.Verification.Ticket,
		Code:       code,
	}
	respData := struct {
		Token string `json:"verificationToken"`
	}{}

	err := t.Client.Execute("POST", "twostepverification", path(
		"v1/users/%d/challenges/%s/verify", nil, l.User.ID, verificationPaths[l.Verification.Type],
	), reqData, &respData)
	if err != nil {
		return nil, err
	}

	return &TwoStepToken{
		User:      l.User.ID,
		Challenge: l.Verification.Ticket,
		Token:     respData.Token,
	}, nil
}
