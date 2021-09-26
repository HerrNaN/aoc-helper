package aoc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

var (
	ErrNoSession = errors.New("no session provided")
)

const (
	sessionCookieName = "session"
	baseURLString     = "https://adventofcode.com"
	inputCacheDir     = ".aoc/input"

	cacheFilePerm os.FileMode = 0755
)

func (h *helper) GetInput(session string) (string, error) {
	cachedInput, err := h.getCachedInput()
	if err == nil {
		return cachedInput, nil
	}

	if session == "" {
		return "", ErrNoSession
	}

	req, err := http.NewRequest(http.MethodGet, h.createGetInputURL(), nil)
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

	h.cacheInput(input)

	return string(input), nil
}

func createSessionCookie(session string) *http.Cookie {
	return &http.Cookie{
		Name:     sessionCookieName,
		Value:    session,
		HttpOnly: true,
	}
}

func (h *helper) createGetInputURL() string {
	return fmt.Sprintf("%s/%d/day/%d/input", baseURLString, h.year, h.day)
}

func (h *helper) getCachedInput() (string, error) {
	cachePath := h.createCachePath()

	stat, err := h.fs.Stat(cachePath)
	if err != nil {
		return "", err
	}

	if stat.IsDir() {
		return "", fmt.Errorf("cache file [%s] is a directory", cachePath)
	}

	input, err := h.fs.ReadFile(cachePath)
	if err != nil {
		return "", err
	}

	return string(input), nil
}

func (h *helper) cacheInput(input []byte) {
	cachePath := h.createCachePath()

	err := h.fs.MkdirAll(path.Dir(cachePath), cacheFilePerm)
	if err != nil {
		return
	}

	err = h.fs.WriteFile(cachePath, input, cacheFilePerm)
	if err != nil {
		return
	}
}

func (h *helper) createCachePath() string {
	return fmt.Sprintf("%s/%s/%d/%02d", h.homeDir, inputCacheDir, h.year, h.day)
}
