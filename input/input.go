package input

import (
	"github.com/gen2brain/dlgs"
)

// Get will prompt the user for files and return a string array of file paths to those files.
func Get() ([]string, *bool, error) {
	files, confirm, err := dlgs.FileMulti("Select files", "")
	if err != nil {
		return nil, nil, err
	}
	return files, &confirm, nil
}
