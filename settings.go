package main

import (
	"github.com/naoina/toml"
	"os"
)

type Settings struct {
	Feeds struct {
		Weather string
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
