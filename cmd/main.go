package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/0xAX/notificator"
	"github.com/NathanReginato/filemaster/file"

	"github.com/NathanReginato/filemaster/activity"
	"github.com/barsanuphe/goexiftool"
	yaml "gopkg.in/yaml.v2"
)

var notify *notificator.Notificator

func main() {

	// Get file path strings from user input
	files, err := file.Get()
	if err != nil {
		panic(err)
	}

	// Get the user activity for the given day
	activity, err := activity.Get()
	if err != nil {
		panic(err)
	}

	fmt.Println(activity)

	for _, v := range files {
		mediaFile := makeMediaFile(v)
		mediaType := getFileType(mediaFile)

		fmt.Println("Media Type: ", mediaType)

		copyMedia(mediaFile, *activity)
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

func makeMediaFile(filePath string) *goexiftool.MediaFile {

	fmt.Println("File to be read: ", filePath)
	mediaFile, err := goexiftool.NewMediaFile(filePath)
	if err != nil {
		panic(err)
	}
	return mediaFile
}

type mediaType int

const (
	video mediaType = 0
	image mediaType = 1
)

func getFileType(mediaFile *goexiftool.MediaFile) mediaType {
	mime, err := mediaFile.Get("MIME Type")
	if err != nil {
		panic(err)
	}
	fmt.Println("MIME Type: ", mime)

	mimeType := strings.Split(mime, "/")[0]
	fmt.Println("Parsed MIME Type: ", mimeType)

	mimeT := image

	if mimeType == "image" {
		mimeT = image
	} else {
		mimeT = video
	}
	return mimeT
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
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(newFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
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
				log.Fatal(err)
			}

			year := date.Year()

			path += strconv.Itoa(year) + "/"

		case "event":
			path += activity + "/"
		case "month-day":

			date, err := mediaFile.GetDate()
			if err != nil {
				log.Fatal(err)
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
