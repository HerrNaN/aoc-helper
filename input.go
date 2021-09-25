package aoc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	ErrNoSession = errors.New("no session provided")
)

const (
	sessionCookieName = "session"
	baseURLString     = "https://adventofcode.com"
)

func (h *Helper) GetInput(session string, day Day, year int) (string, error) {
	if err := day.Validate(); err != nil {
		return "", fmt.Errorf("invalid day: %w", err)
	}

	if session == "" {
		return "", ErrNoSession
	}

	req, err := http.NewRequest(http.MethodGet, createGetInputURL(day, year), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.AddCookie(createSessionCookie(session))

	response, err := h.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform reqest: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return "", errors.New("server responded with non-200 status code")
	}

	input, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(input), nil
}

func createSessionCookie(session string) *http.Cookie {
	return &http.Cookie{
		Name:     sessionCookieName,
		Value:    session,
		HttpOnly: true,
	}
}

func createGetInputURL(day Day, year int) string {
	return fmt.Sprintf("%s/%d/day/%d/input", baseURLString, year, day)
}
