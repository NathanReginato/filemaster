package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/0xAX/notificator"
	"github.com/NathanReginato/filemaster/path"

	"github.com/NathanReginato/filemaster/file"

	"github.com/NathanReginato/filemaster/activity"
	"github.com/barsanuphe/goexiftool"
	yaml "gopkg.in/yaml.v2"
)

var (
	notify *notificator.Notificator
)

func main() {

	// Set up logger for debugging purposes (turn into flag)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Get file path strings from user input
	ps, err := path.Get()
	if err != nil {
		log.Error().Msgf("Failed to retrieve file paths: %v", err)
	}
	log.Debug().Msgf("Retrived %d files from user", len(ps))

	// Get the user activity for the given day
	a, err := activity.Get()
	if err != nil {
		log.Error().Msgf("Failed to retrieve user activity: %v", err)
	}
	log.Debug().Msgf("User activity collected: '%s'", *a)

	// Iterate over files and copy them into folders
	for _, p := range ps {

		log.Debug().Msgf("Reading file: %s", p)

		f, err := file.New(p)
		if err != nil {
			log.Error().Msgf("Failed to create file `%s` from path: %v", p, err)
		}

		t, _ := f.GetType()

		fmt.Println("Media Type: ", t)

		copyMedia(f.Get(), *a)
	}

	notifyFinished()
}

func notifyFinished() {
	notify = notificator.New(notificator.Options{
		DefaultIcon: "icon/default.png",
		AppName:     "Organizer",
	})

	notify.Push("Organizer", "File Organization Complete", "/home/user/icon.png", notificator.UR_CRITICAL)
}

func getDate(mediaFile *goexiftool.MediaFile) {
	date, err := mediaFile.GetDate()
	if err != nil {
		panic(err)
	}
	fmt.Println("Date: ", date)
}

func copyMedia(mediaFile *goexiftool.MediaFile, activity string) {

	fileName := getNewFileName(mediaFile.Filename)
	newFile := buildDirectoryStructure(mediaFile, activity)

	fmt.Println("Filename: ", fileName)

	fmt.Println("New File : ", newFile)

	from, err := os.Open(mediaFile.Filename)

	if err != nil {
		log.Print(err)
	}
	defer from.Close()

	to, err := os.OpenFile(newFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Print(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Print(err)
	}
}

func getNewFileName(filePath string) string {
	path := strings.Split(filePath, "/")
	return path[len(path)-1]
}

type directoryLevel string

const (
	year       directoryLevel = "year"
	event      directoryLevel = "event"
	monthDay   directoryLevel = "month-day"
	cameraType directoryLevel = "camera-type"
	mediatype  directoryLevel = "photo-video"
)

func buildDirectoryStructure(mediaFile *goexiftool.MediaFile, activity string) string {

	config := getConfig()
	structure := config.Structure

	path := "/"

	for _, directoryLevel := range structure {
		switch directoryLevel {
		case "year":

			date, err := mediaFile.GetDate()
			if err != nil {
				log.Print(err)
			}

			year := date.Year()

			path += strconv.Itoa(year) + "/"

		case "event":
			path += activity + "/"
		case "month-day":

			date, err := mediaFile.GetDate()
			if err != nil {
				log.Print(err)
			}

			month := date.Month().String()

			day := date.Day()

			path += month + " " + strconv.Itoa(day) + "/"
		case "camera-type":
			mime, err := mediaFile.Get("MIME Type")
			if err != nil {
				panic(err)
			}

			mimeType := strings.Split(mime, "/")[0]

			if mimeType == "image" {
				camera, err := mediaFile.GetCamera()
				if err != nil {
					panic(err)
				}
				path += camera + "/"
			} else {
				camera, err := mediaFile.Get("Model")
				if err != nil {
					panic(err)
				}
				path += camera + "/"
			}
		case "photo-video":
			mime, err := mediaFile.Get("MIME Type")
			if err != nil {
				panic(err)
			}
			path += mime + "/"
		}
	}

	os.MkdirAll(getRoot()+path, os.ModePerm)

	return getRoot() + path + getNewFileName(mediaFile.Filename)
}

type config struct {
	Root      string   `yaml:"root-directory"`
	Structure []string `yaml:"file-structure"`
	Process   []string `yaml:"process"`
}

func getConfig() config {
	t := config{}

	absPath, _ := filepath.Abs("./config.yaml")

	dat, err := ioutil.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	yamlerr := yaml.Unmarshal([]byte(dat), &t)
	if yamlerr != nil {
		panic(yamlerr)
	}

	return t
}

func getRoot() string {
	config := getConfig()

	return config.Root
}

func getNewFile(fileName string) string {
	root := getRoot()

	return root + "/" + fileName
}
