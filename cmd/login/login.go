package main

import (
	"log"
	"log/slog"
	"time"

	"github.com/sewnie/rbxweb"
	"github.com/sewnie/rbxweb/services/auth"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	c := rbxweb.NewClient()
	c.Logger = slog.Default()

	if err := c.AuthV1.SetCSRFToken(); err != nil {
		log.Fatalln("init csrf token:", err)
	}

	t, err := c.AuthTokenV1.CreateToken()
	if err != nil {
		log.Fatalln("token create:", err)
	}

	for {
		s, err := c.AuthTokenV1.GetTokenStatus(t)
		if err != nil {
			log.Fatalln("token status:", err)
		}

		if s.Status == "Validated" {
			break
		}

		time.Sleep(2 * time.Second)
	}

	_, err = c.AuthV1.CreateLogin(t.Code, t.PrivateKey, auth.Token)
	if err != nil {
		log.Fatalln("login:", err)
	}
}
