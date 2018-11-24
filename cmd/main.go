package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/0xAX/notificator"
	"github.com/gen2brain/dlgs"
	"github.com/xor-gate/goexif2/exif"
	"github.com/xor-gate/goexif2/mknote"
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

	fmt.Println(absPath)

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

	for k, v := range files {
		fmt.Println(k, v)
		// dat, err := ioutil.ReadFile(v)
		// if err != nil {
		// 	panic(err)
		// }
		//fmt.Println(dat)
		f, err := os.Open(v)
		if err != nil {
			log.Fatal(err)
		}
		x, err := exif.Decode(f)
		if err != nil {
			log.Fatal(err)
		}

		exif.RegisterParsers(mknote.All...)
		camModel, err := x.Get(exif.Model)
		if err != nil {
			log.Fatal(err)
		}

		camMake, err := x.Get(exif.Make)
		if err != nil {
			log.Fatal(err)
		}

		tm, _ := x.DateTime()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(camModel)
		fmt.Println(camMake)
		fmt.Println(tm)
	}

	fmt.Println(t)
}
