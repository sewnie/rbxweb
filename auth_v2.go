package rbxweb

import (
	// "crypto/ecdsa"
	// "crypto/elliptic"
	// "crypto/rand"
	// "crypto/sha256"
	// "crypto/x509"
	// "encoding/base64"
	// "errors"
	// "fmt"
	"net/http"
	// "time"
)

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

type loginIntent struct {
	PublicKey    string `json:"clientPublicKey"`
	Epoch        int64  `json:"clientEpochTimestamp"`
	ServerNonce  string `json:"serverNonce"`
	SaiSignature string `json:"saiSignature"`
}

type loginRequest struct {
	CType    string       `json:"ctype"`
	CValue   string       `json:"cvalue"`
	Password string       `json:"password"`
	Intent   *loginIntent `json:"secureAuthenticationIntent,omitempty"`
}

func (a *AuthServiceV2) setCSRFToken() error {
	req, err := a.Client.NewRequest("POST", "auth", "v2/login", nil)
	if err != nil {
		return err
	}

	resp, err := a.Client.BareDo(req)
	if resp.StatusCode == http.StatusForbidden {
		return nil
	}

	return err
}

// CreateLogin logins as the user with the given Token.
//
// If logging in with a username and password, set value and password to a
// username and password.
// If logging in with a Token, set the login type to LoginTypeToken, the value to
// the token's code, and the password to the token's private key.
func (a *AuthServiceV2) CreateLogin(value, password string, login LoginType) (*Login, error) {
	if err := a.Client.csrfRequired(); err != nil {
		return nil, err
	}

	tracker := &http.Cookie{
		Name: "RBXEventTrackerV2",
		// This can definitely be revoked by Roblox if they care to do so.
		// Once that happens, it will mean rbxweb will have to initialize a tracker
		// when required for a request.
		Value: "CreateDate=08/05/2025 14:55:40&rbxid=&browserid=1754423740417001",
	}

	lreq := loginRequest{
		CType:    string(login),
		CValue:   value,
		Password: password,
		// Intent:   *i,
	}

	req, err := a.Client.NewRequest("POST", "auth", "v2/login", lreq)
	if err != nil {
		return nil, err
	}
	req.AddCookie(tracker)

	l := new(Login)
	if _, err := a.Client.Do(req, &l); err != nil {
		return nil, err
	}

	return l, nil
}

// Intent generation from https://github.com/o3dq/roblox-signature, not necessary after
// some reverse engineering.

// func (a *AuthServiceV2) getServerNonce() (string, error) {
// 	var nonce string

// 	err := a.Client.Execute("GET", "apis", "hba-service/v1/getServerNonce", nil, &nonce)
// 	if err != nil {
// 		return "", err
// 	}

// 	return nonce, nil
// }

// func getLoginIntent(nonce string) (*loginIntent, error) {
// 	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
// 	if err != nil {
// 		return nil, fmt.Errorf("key: %w", err)
// 	}

// 	pubBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
// 	if err != nil {
// 		return nil, fmt.Errorf("public key: %w", err)
// 	}
// 	pub := base64.StdEncoding.EncodeToString(pubBytes)

// 	epoch := time.Now().Unix()
// 	data := fmt.Sprintf("%s|%d|%s", pub, epoch, nonce)

// 	hash := sha256.Sum256([]byte(data))
// 	sig, err := ecdsa.SignASN1(rand.Reader, key, hash[:])
// 	if err != nil {
// 		return nil, fmt.Errorf("sign: %w", err)
// 	}
// 	sai := base64.StdEncoding.EncodeToString(sig)

// 	return &loginIntent{
// 		PublicKey:    pub,
// 		Epoch:        epoch,
// 		ServerNonce:  nonce,
// 		SaiSignature: sai,
// 	}, nil
// }
