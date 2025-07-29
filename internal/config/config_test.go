package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadWithEnvVars(t *testing.T) {
	// Set environment variables
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DB_HOST", "test-host")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "test-user")
	os.Setenv("DB_PASSWORD", "test-password")
	os.Setenv("DB_NAME", "test-db")
	os.Setenv("DB_SSLMODE", "require")
	os.Setenv("GRINEX_BASE_URL", "https://test-grinex.io")
	os.Setenv("GRINEX_TIMEOUT", "60s")
	os.Setenv("GRINEX_USER_AGENT", "TestAgent/1.0")
	os.Setenv("LOG_LEVEL", "debug")

	// Clear environment after test
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_SSLMODE")
		os.Unsetenv("GRINEX_BASE_URL")
		os.Unsetenv("GRINEX_TIMEOUT")
		os.Unsetenv("GRINEX_USER_AGENT")
		os.Unsetenv("LOG_LEVEL")
	}()

	cfg := Load()

	assert.Equal(t, "9090", cfg.Server.Port)
	assert.Equal(t, "test-host", cfg.Database.Host)
	assert.Equal(t, 5433, cfg.Database.Port)
	assert.Equal(t, "test-user", cfg.Database.User)
	assert.Equal(t, "test-password", cfg.Database.Password)
	assert.Equal(t, "test-db", cfg.Database.DBName)
	assert.Equal(t, "require", cfg.Database.SSLMode)
	assert.Equal(t, "https://test-grinex.io", cfg.Grinex.BaseURL)
	assert.Equal(t, 60*time.Second, cfg.Grinex.Timeout)
	assert.Equal(t, "TestAgent/1.0", cfg.Grinex.UserAgent)
	assert.Equal(t, "debug", cfg.Logging.Level)
}

func TestGetDSN(t *testing.T) {
	cfg := &Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "password",
			DBName:   "testdb",
			SSLMode:  "disable",
		},
	}

	dsn := cfg.Database.GetDSN()
	expected := "host=localhost port=5432 user=postgres password=password dbname=testdb sslmode=disable"
	assert.Equal(t, expected, dsn)
}
