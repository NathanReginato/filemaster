package file

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/NathanReginato/filemaster/config"
	"github.com/barsanuphe/goexiftool"
)

// File contains a path and a Media File
type File struct {
	mfile  *goexiftool.MediaFile
	path   string
	config *config.Config
}

// New will load a new file into memory from the path provided
func New(c *config.Config, p string) (File, error) {
	// TODO: Account for unitialized File
	f := File{}

	var err error
	f.path = p
	f.mfile, err = goexiftool.NewMediaFile(p)
	f.config = c

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

// Copy takes the file the method is called on, and copies it to the given directory
func (f *File) Copy(dest string) error {

	p := f.GetPath()

	from, err := os.Open(p)
	if err != nil {
		return fmt.Errorf("unable to open file %s: %v", p, err)
	}
	defer from.Close()

	to, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("unable to create file %s: %v", p, err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return fmt.Errorf("unable to create copy data %s: %v", p, err)
	}

	return nil
}

// GetDestination returns the calcualted path of the file based on the YAML configuration
func (f *File) GetDestination(activity *string) (*string, error) {

	fi := f.Get()
	s := f.config.GetStructure()
	w := f.config.GetWorkspace()

	p := []string{}

	if len(s) == 0 {
		return nil, fmt.Errorf("file-structure array has no elements, please check the config")
	}

	for _, d := range s {
		switch d {
		case "year":

			d, err := fi.GetDate()
			if err != nil {
				log.Print(err)
			}

			year := d.Year()

			a := strconv.Itoa(year) + "/"

			p = append(p, a)

		case "event":
			p = append(p, *activity)
		case "month-day":

			d, err := fi.GetDate()
			if err != nil {
				log.Print(err)
			}

			month := d.Month().String()
			day := strconv.Itoa(d.Day())
			date := month + " " + day

			p = append(p, date)
		case "camera-type":
			t, err := f.GetType()
			if err != nil {
				panic(err)
			}

			if *t == "image" {
				camera, err := fi.GetCamera()
				if err != nil {
					panic(err)
				}
				p = append(p, camera)
			} else {
				camera, err := fi.Get("Model")
				if err != nil {
					panic(err)
				}
				p = append(p, camera)
			}
		case "photo-video":
			mime, err := fi.Get("MIME Type")
			if err != nil {
				panic(err)
			}

			t := strings.Split(mime, "/")

			kind := t[0]
			fileType := t[1]

			p = append(p, kind)
			p = append(p, fileType)
		}
	}

	aliasedPath := f.replaceAliases(p)

	path := string(os.PathSeparator) + strings.Join(aliasedPath, string(os.PathSeparator))
	ab := w + path
	np := ab + string(os.PathSeparator) + f.GetName()

	os.MkdirAll(ab, os.ModePerm)

	return &np, nil
}

func (f *File) replaceAliases(p []string) []string {
	a := f.config.GetAlias()
	for i, segment := range p {
		if val, ok := a[segment]; ok {
			p[i] = val
		}
	}

	return p
}
