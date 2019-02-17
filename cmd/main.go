package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/0xAX/notificator"
	"github.com/NathanReginato/filemaster/path"

	"github.com/NathanReginato/filemaster/file"

	"github.com/NathanReginato/filemaster/config"

	"github.com/NathanReginato/filemaster/activity"
	"github.com/barsanuphe/goexiftool"
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

		fmt.Println("Media Type: ", *t)

		copyMedia(f, *a)
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

func copyMedia(f file.File, activity string) {

	newFile := buildDirectoryStructure(f, activity)

	from, err := os.Open(f.GetPath())

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

type directoryLevel string

const (
	year       directoryLevel = "year"
	event      directoryLevel = "event"
	monthDay   directoryLevel = "month-day"
	cameraType directoryLevel = "camera-type"
	mediatype  directoryLevel = "photo-video"
)

func buildDirectoryStructure(f file.File, activity string) string {

	c, err := config.New()
	if err != nil {
		log.Error().Msgf("loading config failed: %v", err)
	}

	structure := c.GetStructure()

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
			path += activity + "/"
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

	os.MkdirAll(c.GetWorkspace()+path, os.ModePerm)

	return c.GetWorkspace() + path + f.GetName()
}
