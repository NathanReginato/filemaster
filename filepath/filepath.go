package filepath

import (
	"github.com/gen2brain/dlgs"
)

// Get will prompt the user for files and return a string array of file paths to those files.
func Get() ([]string, error) {
	files, _, err := dlgs.FileMulti("Select files", "")
	if err != nil {
		return nil, err
	}
	return files, nil
}
