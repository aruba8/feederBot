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
		assert.Equal(t, "prod", s.Environment)

	})

	t.Run("Test correct connection string", func(t *testing.T) {
		s := GetSettings()
		db := s.Database()
		assert.Equal(t, "mongodb://localhost:27017", db.ConnectionString())
	})

}
