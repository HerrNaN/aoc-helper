package aoc

import "fmt"

type Day int

func (d Day) Validate() error {
	if d < 1 || d > 25 {
		return fmt.Errorf("a days value must be between 1 and 25 inclusive: %d", d)
	}

	return nil
}
