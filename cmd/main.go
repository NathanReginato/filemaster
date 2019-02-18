package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/NathanReginato/filemaster/activity"
	"github.com/NathanReginato/filemaster/config"
	"github.com/NathanReginato/filemaster/file"
	"github.com/NathanReginato/filemaster/input"
	"github.com/NathanReginato/filemaster/notify"
)

var (
	conf *config.Config
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
	i, err := input.Get()
	if err != nil {
		log.Error().Msgf("failed to retrieve file paths: %v", err)
	}
	log.Debug().Msgf("retrived %d files from user", len(i))

	var a *string
	if conf.HasDir("event") {
		// Get the user activity for the given day
		var err error
		a, err = activity.Get()
		if err != nil {
			log.Error().Msgf("failed to retrieve user activity: %v", err)
		}
		log.Debug().Msgf("user activity collected: '%s'", *a)
	}

	// Iterate over files and copy them into folders
	for _, p := range i {

		log.Debug().Msgf("reading file: %s", p)

		f, err := file.New(conf, p)
		if err != nil {
			log.Error().Msgf("failed to create file `%s` from path: %v", p, err)
		}

		directory, err := f.GetDestination(a)
		err = f.Copy(*directory)
		if err != nil {
			log.Error().Msgf("failed to move file to directory: `%v`", err)
		}
	}

	notify.Finished()
}
