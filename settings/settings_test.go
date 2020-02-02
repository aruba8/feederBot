package settings

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSettings(t *testing.T) {
	os.Setenv("SETTINGS_FILE_PATH", "../settings.toml")

	t.Run("Test correct env set if LAMBDA_ENVIRON is local", func(t *testing.T) {
		os.Setenv("LAMBDA_ENVIRON", "local")
		s := GetSettings()
		assert.Equal(t, "local", s.Environment)
	})

	t.Run("Test correct env set if LAMBDA_ENVIRON is not set", func(t *testing.T) {
		s := GetSettings()
		assert.Equal(t, "local", s.Environment)
	})

	t.Run("Test correct env LAMBDA_ENVIRON is prod", func(t *testing.T) {
		os.Setenv("LAMBDA_ENVIRON", "prod")
		s := GetSettings()
		os.Setenv("LAMBDA_ENVIRON", "")
		db := s.Database()
		assert.Equal(t, "prod", s.Environment)
		assert.Equal(t, "mongodb://cluster0-txu1v.mongodb.net:27017", db.ConnectionString())

	})

	t.Run("Test correct connection string LAMBDA_ENVIRON is prod", func(t *testing.T) {
		os.Setenv("LAMBDA_ENVIRON", "prod")
		s := GetSettings()
		db := s.Database()
		os.Setenv("LAMBDA_ENVIRON", "")
		assert.Equal(t, "prod", s.Environment)
		assert.Equal(t, "mongodb://cluster0-txu1v.mongodb.net:27017", db.ConnectionString())

	})

	t.Run("Test correct connection string", func(t *testing.T) {
		s := GetSettings()
		db := s.Database()
		assert.Equal(t, "mongodb://localhost:27017", db.ConnectionString())
	})
}

func TestGetSettings(t *testing.T) {
	t.Run("Test settings file wrong", func(t *testing.T) {
		os.Setenv("SETTINGS_FILE_PATH", "../settings")
		assert.Panics(t, func() {
			GetSettings()
		})
	})
}
