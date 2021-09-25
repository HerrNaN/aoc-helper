package aoc

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

func TestHelper_GetInput_WithValidInputSucceeds(t *testing.T) {
	defer gock.Off()

	session := "some cookie value"
	input := "some input"

	gock.New("https://adventofcode.com").
		Get("/2021/day/13/input").
		MatchHeader("Cookie", session).
		Reply(200).
		JSON(input)

	h := NewHelper()
	actualInput, err := h.GetInput(session, 13, 2021)

	require.NoError(t, err)
	require.Equal(t, input, actualInput)
}

func TestHelper_GetInput_WithoutSessionFails(t *testing.T) {
	h := NewHelper()
	_, err := h.GetInput("", 13, 2021)

	require.Error(t, err)
	require.ErrorIs(t, err, ErrNoSession)
}

func TestHelper_GetInput_WithInvalidSessionFails(t *testing.T) {
	defer gock.Off()

	invalidSession := "invalid session"

	gock.New("https://adventofcode.com").
		Get("/2021/day/13/input").
		MatchHeader("Cookie", invalidSession).
		Reply(500)

	h := NewHelper()
	_, err := h.GetInput(invalidSession, 13, 2021)

	require.Error(t, err)
}
