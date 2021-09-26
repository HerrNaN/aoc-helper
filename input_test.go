package aoc

import (
	"errors"
	"net/http"
	"testing"

	"github.com/spf13/afero"

	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

func TestHelper_GetInput_ShouldFailWithErrNoSessionWithoutSession(t *testing.T) {
	fakeFS := afero.Afero{Fs: new(afero.MemMapFs)}

	h := &helper{
		client: http.DefaultClient,
		day:    13,
		year:   2021,
		fs:     fakeFS,
	}

	_, err := h.GetInput("")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNoSession)
}

func TestHelper_GetInput_ShouldFailWithInvalidSession(t *testing.T) {
	defer gock.Off()

	invalidSession := "invalid session"

	gock.New("https://adventofcode.com").
		Get("/2021/day/13/input").
		MatchHeader("Cookie", invalidSession).
		Reply(500)

	fakeFS := afero.Afero{Fs: new(afero.MemMapFs)}

	h := &helper{
		client: http.DefaultClient,
		day:    13,
		year:   2021,
		fs:     fakeFS,
	}

	_, err := h.GetInput(invalidSession)
	require.Error(t, err)
}

func TestHelper_GetInput_ShouldUseCachedInputWhenItExists(t *testing.T) {
	defer gock.Off()

	expectedInput := "test"

	gock.New("https://adventofcode.com").
		Get("/2021/day/13/input").
		ReplyError(errors.New(""))

	fakeFS := afero.Afero{Fs: new(afero.MemMapFs)}
	fakeFS.MkdirAll("/home/test/.aoc/input/2021", 0755)
	fakeFS.WriteFile("/home/test/.aoc/input/2021/13", []byte(expectedInput), 0755)

	h := &helper{
		client:  http.DefaultClient,
		day:     13,
		year:    2021,
		fs:      fakeFS,
		homeDir: "/home/test",
	}

	actualInput, err := h.GetInput("")

	require.NoError(t, err)
	require.Equal(t, expectedInput, actualInput)
}

func TestHelper_GetInput_ShouldDownloadInputWhenCacheDoesntExist(t *testing.T) {
	defer gock.Off()

	expectedInput := "test"
	session := "valid session"

	gock.New("https://adventofcode.com").
		Get("/2021/day/13/input").
		MatchHeader("Cookie", session).
		Reply(200).
		JSON(expectedInput)

	fakeFS := afero.Afero{Fs: new(afero.MemMapFs)}

	h := &helper{
		client: http.DefaultClient,
		day:    13,
		year:   2021,
		fs:     fakeFS,
	}

	actualInput, err := h.GetInput(session)

	require.NoError(t, err)
	require.Equal(t, expectedInput, actualInput)
}

func TestHelper_GetInput_ShouldCacheDownloadedInput(t *testing.T) {
	defer gock.Off()

	expectedInput := "test"
	session := "valid session"

	gock.New("https://adventofcode.com").
		Get("/2021/day/13/input").
		MatchHeader("Cookie", session).
		Reply(200).
		JSON(expectedInput)

	fakeFS := afero.Afero{Fs: new(afero.MemMapFs)}

	h := &helper{
		client:  http.DefaultClient,
		day:     13,
		year:    2021,
		fs:      fakeFS,
		homeDir: "/home/test",
	}

	_, err := h.GetInput(session)
	require.NoError(t, err)

	cacheFile := h.homeDir + "/.aoc/input/2021/13"

	didCache, err := h.fs.FileContainsBytes(cacheFile, []byte(expectedInput))
	require.NoError(t, err)
	require.True(t, didCache)

}
