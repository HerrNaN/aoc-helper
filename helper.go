package aoc

import (
	"net/http"
	"os"

	"github.com/spf13/afero"
)

type Day int
type Year int

type helper struct {
	client  *http.Client
	day     Day
	year    Year
	fs      afero.Afero
	homeDir string
}

func NewHelper(day Day, year Year) (*helper, error) {
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
