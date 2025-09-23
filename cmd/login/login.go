package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/sewnie/rbxweb"
)

type debugTransport struct {
	underlying http.RoundTripper
}

func (t *debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	slog.Debug("HTTP request",
		"method", req.Method,
		"url", req.URL.String(),
	)

	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewReader(body))
		slog.Debug("HTTP request body", "body", string(body))
	}

	resp, err := t.underlying.RoundTrip(req)
	if err != nil {
		slog.Debug("HTTP request failed",
			"error", err,
		)
		return nil, err
	}

	slog.Debug("HTTP response",
		"status", resp.StatusCode,
	)
	return resp, nil
}

func main() {
	c := rbxweb.NewClient()
	slog.SetLogLoggerLevel(slog.LevelDebug)
	c.Client.Transport = &debugTransport{
		underlying: http.DefaultTransport,
	}

	if len(os.Args) == 3 {
		l, err := c.AuthV2.CreateLogin(os.Args[1], os.Args[2], rbxweb.LoginTypeUsername)
		if err != nil {
			log.Fatalln("login create:", err)
		}

		log.Printf("%+v", l)

		if v := l.Verification; v != nil {
			code := readPrompt("Code:")

			tsToken, err := c.VerificationServiceV1.CreateToken(l, code)
			if err != nil {
				log.Fatalln("verify:", err)
			}
			log.Printf("%+v", tsToken)

			tsLogin, err := c.AuthV3.CreateTwoStepLogin(false, tsToken)
			if err != nil {
				log.Fatalln("verify login:", err)
			}
			log.Printf("%+v", tsLogin)
		}

		os.Exit(0)
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
		log.Println(s)

		if s.Status == "Validated" {
			break
		}

		time.Sleep(4 * time.Second)
	}

	log.Fatal(c.AuthV2.CreateLogin(t.Code, t.PrivateKey, rbxweb.LoginTypeToken))
}

func readPrompt(format string, a ...any) string {
	log.Printf(format, a...)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		log.Fatalln("scanner:", err)
	}

	data := scanner.Text()
	if len(data) == 0 {
		log.Fatal("zero length input")
	}

	return data
}
