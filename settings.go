package main

import (
	"fmt"
	"github.com/naoina/toml"
	"os"
)

type Settings struct {
	Feeds struct {
		Weather string
	}
	Database struct {
		Name       string
		Collection string
		URI        string
		Port       int
	}
}

func getSettings() Settings {
	f, err := os.Open("settings.toml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var settings Settings
	if err := toml.NewDecoder(f).Decode(&settings); err != nil {
		panic(err)
	}
	return settings
}

func (s *Settings) getDbUri() string {
	uriTmpl := "mongodb://%s:%d"
	uri := fmt.Sprintf(uriTmpl, s.Database.URI, s.Database.Port)
	return uri
}
