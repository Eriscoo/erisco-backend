package turnstile

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type verifyResponse struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

type Service struct {
	secretKey string
	client    *http.Client
}

func New(secretKey string) *Service {
	return &Service{
		secretKey: secretKey,
		client:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *Service) Verify(token string) (bool, error) {
	resp, err := s.client.PostForm(
		"https://challenges.cloudflare.com/turnstile/v0/siteverify",
		url.Values{
			"secret":   {s.secretKey},
			"response": {token},
		},
	)
	if err != nil {
		return false, fmt.Errorf("turnstile verify: %w", err)
	}
	defer resp.Body.Close()

	var v verifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return false, fmt.Errorf("turnstile decode: %w", err)
	}

	return v.Success, nil
}
