package token

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
