package activity

import (
	"errors"

	"github.com/gen2brain/dlgs"
)

// Get returns a user input
func Get() (*string, error) {
	a, s, err := dlgs.Entry("Activity", "What did you do this day?", "Took product shots")
	if err != nil {
		return nil, err
	}

	if s == false {
		return nil, errors.New("activity was captured from user with out error, but was still successful")
	}

	return &a, nil
}
