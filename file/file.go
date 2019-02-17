package file

import (
	"strings"

	"github.com/barsanuphe/goexiftool"
)

// File contains a path and a Media File
type File struct {
	mfile *goexiftool.MediaFile
}

// New will load a new file into memory from the path provided
func New(p string) (File, error) {
	f := File{}

	var err error
	f.mfile, err = goexiftool.NewMediaFile(p)

	return f, err
}

// Get will load a new file into memory from the path provided
func (f *File) Get() *goexiftool.MediaFile {
	return f.mfile
}

// GetType will return the MIME type of the passed Media File
func (f *File) GetType() (*string, error) {

	m, err := f.mfile.Get("MIME Type")
	if err != nil {
		return nil, err
	}

	t := strings.Split(m, "/")[0]

	return &t, nil
}
