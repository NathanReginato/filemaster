package file

import (
	"strings"

	"github.com/barsanuphe/goexiftool"
)

// File contains a path and a Media File
type File struct {
	mfile *goexiftool.MediaFile
	path  string
}

// New will load a new file into memory from the path provided
func New(p string) (File, error) {
	// TODO: Account for unitialized File
	f := File{}

	var err error
	f.path = p
	f.mfile, err = goexiftool.NewMediaFile(p)

	return f, err
}

// Get will load a new file into memory from the path provided
func (f *File) Get() *goexiftool.MediaFile {
	// TODO: Account for unitialized File
	return f.mfile
}

// GetType will return the MIME type of the passed Media File
func (f *File) GetType() (*string, error) {
	// TODO: Account for unitialized File

	m, err := f.mfile.Get("MIME Type")
	if err != nil {
		return nil, err
	}

	t := strings.Split(m, "/")[0]

	return &t, nil
}

// GetName will take the file path of the file and return it's name
func (f *File) GetName() string {
	// TODO: Account for unitialized File
	path := strings.Split(f.path, "/")
	return path[len(path)-1]
}

// GetPath will return the full path to the file
func (f *File) GetPath() string {
	// TODO: Account for unitialized File
	return f.path
}
