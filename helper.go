package aoc

import (
	"net/http"
	"os"

	"github.com/spf13/afero"
)

type Day int
type Year int

type Helper struct {
	client  *http.Client
	day     Day
	year    Year
	fs      afero.Afero
	homeDir string
}

func NewHelper(day Day, year Year) (*Helper, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return &Helper{
		day:    day,
		year:   year,
		client: http.DefaultClient,
		fs: afero.Afero{
			Fs: afero.NewOsFs(),
		},
		homeDir: home,
	}, nil
}
