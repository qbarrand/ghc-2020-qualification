package pkg

import (
	"fmt"
	"os"
)

type input struct {
}

func ParseInput(filename string) (*input, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	in := input{}

	lineno := 1

	if _, err := fmt.Scanf("some-format", &in); err != nil {
		return nil, fmt.Errorf("could not read the header: %v", err)
	}

	lineno++

	for { // TODO: add a condition
		if _, err := fmt.Scanf("some-format", &in); err != nil {
			return nil, fmt.Errorf("could not read line %d: %v", lineno, err)
		}

		lineno++
	}

	return &in, nil
}
