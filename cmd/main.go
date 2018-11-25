package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/0xAX/notificator"
	"github.com/barsanuphe/goexiftool"
	"github.com/gen2brain/dlgs"
	yaml "gopkg.in/yaml.v2"
)

var notify *notificator.Notificator

type config struct {
	Root      string   `yaml:"root-directory"`
	Structure []string `yaml:"file-structure"`
	Process   []string `yaml:"process"`
}

func main() {

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

	notify = notificator.New(notificator.Options{
		DefaultIcon: "icon/default.png",
		AppName:     "My test App",
	})

	notify.Push("title", "text", "/home/user/icon.png", notificator.UR_CRITICAL)

	files, _, err := dlgs.FileMulti("Select files", "")
	if err != nil {
		panic(err)
	}

	printFileMetaData(files)
}

type mediaType int

const (
	video mediaType = 0
	image mediaType = 1
)

func printFileMetaData(files []string) {
	for _, v := range files {

		mediaFile := makeMediaFile(v)
		mediaType := getFileType(mediaFile)

		fmt.Println("Media Type: ", mediaType)

		if mediaType == image {
			handleImage(mediaFile)
			getDate(mediaFile)
		} else {
			handleVideo(mediaFile)
			getDate(mediaFile)
		}

	}
}

func makeMediaFile(filePath string) *goexiftool.MediaFile {

	fmt.Println("File to be read: ", filePath)
	mediaFile, err := goexiftool.NewMediaFile(filePath)
	if err != nil {
		panic(err)
	}
	return mediaFile
}

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

func handleImage(mediaFile *goexiftool.MediaFile) {
	camera, err := mediaFile.GetCamera()
	if err != nil {
		panic(err)
	}
	fmt.Println("Make: ", camera)
}

func handleVideo(mediaFile *goexiftool.MediaFile) {
	camera, err := mediaFile.Get("Model")
	if err != nil {
		panic(err)
	}
	fmt.Println("Make: ", camera)
}

func getDate(mediaFile *goexiftool.MediaFile) {
	date, err := mediaFile.GetDate()
	if err != nil {
		panic(err)
	}
	fmt.Println("Date: ", date)
}
