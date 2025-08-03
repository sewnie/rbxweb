package main

import (
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

	log.Println("Retrieved CSRF token:", c.Token)

	t, err := c.AuthTokenV1.CreateToken()
	if err != nil {
		log.Fatalln("token create:", err)
	}
	log.Println("Retrieved token:", t.Code, t.PrivateKey)

	for {
		s, err := c.AuthTokenV1.GetTokenStatus(t)
		if err != nil {
			log.Fatalln("token status:", err)
		}

		log.Printf("%+v", s)

		if s.Status == "Validated" {
			break
		}

		time.Sleep(2 * time.Second)
	}

	r, err := c.AuthV1.CreateLogin(t.Code, t.PrivateKey, auth.Token)
	if err != nil {
		log.Fatalln("login:", err)
	}

	log.Printf("%+v", r)
	log.Println("ROBLOSECURITY:", c.Security)
}
