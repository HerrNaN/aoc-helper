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
	ErrNon200Response = errors.New("server responded with non-200 status code")
)

const (
	baseURLString     = "https://adventofcode.com"
	sessionCookieName = "session"

	sessionFile               = ".aoc/session"
	inputCacheDir             = ".aoc/input"
	cacheFilePerm os.FileMode = 0755
)

// GetSession retreives the input for the day and year specified in the helper.
func (h *Helper) GetInput() (string, error) {
	cachedInput, err := h.getCachedInput()
	if err == nil {
		return cachedInput, nil
	}

	input, err := h.downloadInput()
	if err != nil {
		return "", fmt.Errorf("couldn't download input: %w", err)
	}

	h.cacheInput(input)

	return string(input), nil
}

func (h *Helper) downloadInput() ([]byte, error) {
	session, err := h.getSession()
	if err != nil {
		return nil, fmt.Errorf("couldn't get session: %w", err)
	}

	return h.downloadInputUsingSession(session)
}

func (h *Helper) downloadInputUsingSession(session string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, h.createGetInputURL(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.AddCookie(createSessionCookie(session))

	response, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform reqest: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, ErrNon200Response
	}

	input, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return input, nil
}

func createSessionCookie(session string) *http.Cookie {
	return &http.Cookie{
		Name:     sessionCookieName,
		Value:    session,
		HttpOnly: true,
	}
}

func (h *Helper) getSession() (string, error) {
	sessionPath := h.sessionPath()

	stat, err := h.fs.Stat(sessionPath)
	if err != nil {
		return "", fmt.Errorf("couldn't get file info for '%s': %w", sessionPath, err)
	}

	if stat.IsDir() {
		return "", fmt.Errorf("session file '%s' is a directory", sessionPath)
	}

	session, err := h.fs.ReadFile(sessionPath)
	if err != nil {
		return "", fmt.Errorf("couldn't read session file '%s': %w", sessionPath, err)
	}

	return string(session), nil

}

func (h *Helper) sessionPath() string {
	return fmt.Sprintf("%s/%s", h.homeDir, sessionFile)
}

func (h *Helper) createGetInputURL() string {
	return fmt.Sprintf("%s/%d/day/%d/input", baseURLString, h.year, h.day)
}

func (h *Helper) getCachedInput() (string, error) {
	cachePath := h.getCachePath()

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

func (h *Helper) cacheInput(input []byte) {
	cachePath := h.getCachePath()

	err := h.fs.MkdirAll(path.Dir(cachePath), cacheFilePerm)
	if err != nil {
		return
	}

	err = h.fs.WriteFile(cachePath, input, cacheFilePerm)
	if err != nil {
		return
	}
}

func (h *Helper) getCachePath() string {
	return fmt.Sprintf("%s/%s/%d/%02d", h.homeDir, inputCacheDir, h.year, h.day)
}
