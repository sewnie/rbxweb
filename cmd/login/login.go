package main

import (
	"os"
	"log"
	"log/slog"
	"time"

	"github.com/sewnie/rbxweb"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	c := rbxweb.NewClient()
	c.Logger = slog.Default()

	if len(os.Args) == 3 {
		log.Fatal(c.AuthV2.CreateLogin(os.Args[1], os.Args[2], rbxweb.LoginTypeUsername))
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

		time.Sleep(4 * time.Second)
	}

	log.Fatal(c.AuthV2.CreateLogin(t.Code, t.PrivateKey, rbxweb.LoginTypeToken))
}
