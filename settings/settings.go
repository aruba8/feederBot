package settings

import (
	"fmt"
	"github.com/naoina/toml"
	"os"
	"path/filepath"
)

type Settings struct {
	Feeds struct {
		Weather string
	}
	Databases   map[string]Database
	Environment string
	BotSettings BotSettings
}

type Database struct {
	Username     string
	Collection   string
	Host         string
	Port         int
	Password     string
	DatabaseName string
}

type BotSettings struct {
	Token string
}

func GetSettings() Settings {
	settingsPath := os.Getenv("SETTINGS_FILE_PATH")
	lambdaEnv := os.Getenv("LAMBDA_ENVIRON")
	f, err := os.Open(filepath.Clean(settingsPath))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var settings Settings
	if len(lambdaEnv) == 0 {
		lambdaEnv = "local"
	}
	if err := toml.NewDecoder(f).Decode(&settings); err != nil {
		panic(err)
	}
	settings.Environment = lambdaEnv
	return settings
}

func (d *Database) ConnectionString() string {
	password := os.Getenv("DB_PASSWORD")
	lambdaEnv := os.Getenv("LAMBDA_ENVIRON")
	connectionString := ""
	if lambdaEnv == "local" || len(lambdaEnv) == 0 {
		connectionString = fmt.Sprintf("mongodb://%s:%d", d.Host, d.Port)
	} else {
		connectionString = fmt.Sprintf("mongodb+srv://%s:%s@%s", d.Username, password, d.Host)
	}
	return connectionString
}

func (s *Settings) Database() Database {
	database := s.Databases[s.Environment]
	database.Password = os.Getenv("DB_PASSWORD")
	return database
}

func (s *Settings) Bot() BotSettings {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	return BotSettings{token}
}
