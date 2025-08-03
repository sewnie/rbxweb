// clientsettings selectively implements the 'auth' Roblox Web API.
package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/sewnie/rbxweb/internal/api"
)

// GamesServiceV1 handles the 'auth/v2' Roblox Web API.
type AuthServiceV1 api.Service

// SetCSRFToken calls v2/login, in hopes of returning a x-csrf-token, handled and set by client.
//
// If the request returned is not 403 Forbidden, an error will be returned.
func (a *AuthServiceV1) SetCSRFToken() error {
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
// If logging in with a token.Token, set the login type to Token, the value to
// the token's code, and the password to the token's private key.
//
// Implementation from https://github.com/o3dq/roblox-signature
//
// Requires a CSRF Token to be set, see SetCSRFToken
func (a *AuthServiceV1) CreateLogin(value, password string, login LoginType) (*Login, error) {
	var r Login

	i, err := a.getLoginIntent()
	if err != nil {
		return nil, err
	}

	req := loginRequest{
		CType:    string(login),
		CValue:   value,
		Password: password,
		Intent:   *i,
	}

	err = a.Client.Execute("POST", "auth", "v2/login", req, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (a *AuthServiceV1) getServerNonce() (string, error) {
	var nonce string

	err := a.Client.Execute("GET", "apis", "hba-service/v1/getServerNonce", nil, &nonce)
	if err != nil {
		return "", err
	}

	return nonce, nil
}

func (a *AuthServiceV1) getLoginIntent() (*loginIntent, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("key: %w", err)
	}

	nonce, err := a.getServerNonce()
	if err != nil {
		return nil, fmt.Errorf("nonce: %w", err)
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("public key: %w", err)
	}
	pub := base64.StdEncoding.EncodeToString(pubBytes)

	epoch := time.Now().Unix()
	data := fmt.Sprintf("%s:%d:%s", pub, epoch, nonce)

	hash := sha256.Sum256([]byte(data))
	sig, err := ecdsa.SignASN1(rand.Reader, key, hash[:])
	if err != nil {
		return nil, fmt.Errorf("sign: %w", err)
	}
	sai := base64.StdEncoding.EncodeToString(sig)

	return &loginIntent{
		PublicKey:    pub,
		Epoch:        epoch,
		ServerNonce:  nonce,
		SaiSignature: sai,
	}, nil
}
