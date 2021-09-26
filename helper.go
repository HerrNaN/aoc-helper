package aoc

import (
	"net/http"
	"os"

	"github.com/spf13/afero"
)

type helper struct {
	client  *http.Client
	day     Day
	year    int
	fs      afero.Afero
	homeDir string
}

func NewHelper(day Day, year int) (*helper, error) {
	if err := day.Validate(); err != nil {
		return nil, err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return &helper{
		day:    day,
		year:   year,
		client: http.DefaultClient,
		fs: afero.Afero{
			Fs: afero.NewOsFs(),
		},
		homeDir: home,
	}, nil
}
