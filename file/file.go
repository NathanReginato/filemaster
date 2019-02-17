package file

import (
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

type directoryLevel string

const (
	year       directoryLevel = "year"
	event      directoryLevel = "event"
	monthDay   directoryLevel = "month-day"
	cameraType directoryLevel = "camera-type"
	mediatype  directoryLevel = "photo-video"
)

func (f *File) Copy(dest string) error {

	from, err := os.Open(f.GetPath())

	if err != nil {
		log.Print(err)
	}
	defer from.Close()

	to, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Print(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Print(err)
	}

	return nil
}

func (f *File) GetDestination(activity *string) string {

	structure := f.config.GetStructure()

	path := "/"

	for _, directoryLevel := range structure {
		switch directoryLevel {
		case "year":

			date, err := f.Get().GetDate()
			if err != nil {
				log.Print(err)
			}

			year := date.Year()

			path += strconv.Itoa(year) + "/"

		case "event":
			path += *activity + "/"
		case "month-day":

			date, err := f.Get().GetDate()
			if err != nil {
				log.Print(err)
			}

			month := date.Month().String()

			day := date.Day()

			path += month + " " + strconv.Itoa(day) + "/"
		case "camera-type":
			mime, err := f.Get().Get("MIME Type")
			if err != nil {
				panic(err)
			}

			mimeType := strings.Split(mime, "/")[0]

			if mimeType == "image" {
				camera, err := f.Get().GetCamera()
				if err != nil {
					panic(err)
				}
				path += camera + "/"
			} else {
				camera, err := f.Get().Get("Model")
				if err != nil {
					panic(err)
				}
				path += camera + "/"
			}
		case "photo-video":
			mime, err := f.Get().Get("MIME Type")
			if err != nil {
				panic(err)
			}
			path += mime + "/"
		}
	}

	os.MkdirAll(f.config.GetWorkspace()+path, os.ModePerm)

	return f.config.GetWorkspace() + path + f.GetName()
}
