package main

import (
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
)

var (
	notify *notificator.Notificator
	conf   *config.Config
)

func main() {

	// Set up logger for debugging purposes (turn into flag)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Get configuration for app
	var err error
	conf, err = config.New()
	if err != nil {
		log.Error().Msgf("loading config failed: %v", err)
	}
	log.Debug().Msgf("loaded config file")

	// Get file path strings from user input
	ps, err := path.Get()
	if err != nil {
		log.Error().Msgf("failed to retrieve file paths: %v", err)
	}
	log.Debug().Msgf("retrived %d files from user", len(ps))

	// Get the user activity for the given day
	a, err := activity.Get()
	if err != nil {
		log.Error().Msgf("failed to retrieve user activity: %v", err)
	}
	log.Debug().Msgf("user activity collected: '%s'", *a)

	// Iterate over files and copy them into folders
	for _, p := range ps {

		log.Debug().Msgf("reading file: %s", p)

		f, err := file.New(p)
		if err != nil {
			log.Error().Msgf("failed to create file `%s` from path: %v", p, err)
		}

		copyMedia(f, *a)
	}

	notifyFinished()
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

	structure := conf.GetStructure()

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

	os.MkdirAll(conf.GetWorkspace()+path, os.ModePerm)

	return conf.GetWorkspace() + path + f.GetName()
}

func notifyFinished() {
	notify = notificator.New(notificator.Options{
		DefaultIcon: "icon/default.png",
		AppName:     "Organizer",
	})

	notify.Push("Organizer", "File Organization Complete", "/home/user/icon.png", notificator.UR_CRITICAL)
}
