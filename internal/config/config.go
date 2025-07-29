package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Grinex   GrinexConfig   `mapstructure:"grinex"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type GrinexConfig struct {
	BaseURL   string        `mapstructure:"base_url"`
	Timeout   time.Duration `mapstructure:"timeout"`
	UserAgent string        `mapstructure:"user_agent"`
}

type LoggingConfig struct {
	Level string `mapstructure:"level"`
}

// Load loads configuration from environment variables and command line flags
func Load() *Config {

	setDefaults()

	loadFromEnv()

	loadFromFlags()

	cfg := &Config{
		Server: ServerConfig{
			Port: getString("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getString("DB_HOST", "localhost"),
			Port:     getInt("DB_PORT", 5460),
			User:     getString("DB_USER", "db_admin"),
			Password: getString("DB_PASSWORD", "3Qv@e8U0ImT"),
			DBName:   getString("DB_NAME", "grinex_rates"),
			SSLMode:  getString("DB_SSLMODE", "disable"),
		},
		Grinex: GrinexConfig{
			BaseURL:   getString("GRINEX_BASE_URL", "https://grinex.io"),
			Timeout:   getDuration("GRINEX_TIMEOUT", 30*time.Second),
			UserAgent: getString("GRINEX_USER_AGENT", "GrinexRateService/1.0"),
		},
		Logging: LoggingConfig{
			Level: getString("LOG_LEVEL", "info"),
		},
	}

	return cfg
}

// GetDSN returns the PostgreSQL connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func setDefaults() {
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.dbname", "grinex_rates")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("grinex.base_url", "https://grinex.io")
	viper.SetDefault("grinex.timeout", "30s")
	viper.SetDefault("grinex.user_agent", "GrinexRateService/1.0")
	viper.SetDefault("logging.level", "info")
}

func loadFromEnv() {
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()
}

func loadFromFlags() {
	port := flag.String("port", "", "Server port")
	dbHost := flag.String("db-host", "", "Database host")
	dbPort := flag.Int("db-port", 0, "Database port")
	dbUser := flag.String("db-user", "", "Database user")
	dbPassword := flag.String("db-password", "", "Database password")
	dbName := flag.String("db-name", "", "Database name")
	dbSSLMode := flag.String("db-sslmode", "", "Database SSL mode")
	grinexBaseURL := flag.String("grinex-base-url", "", "Grinex API base URL")
	grinexTimeout := flag.String("grinex-timeout", "", "Grinex API timeout")
	logLevel := flag.String("log-level", "", "Log level")

	flag.Parse()

	if *port != "" {
		viper.Set("server.port", *port)
	}
	if *dbHost != "" {
		viper.Set("database.host", *dbHost)
	}
	if *dbPort != 0 {
		viper.Set("database.port", *dbPort)
	}
	if *dbUser != "" {
		viper.Set("database.user", *dbUser)
	}
	if *dbPassword != "" {
		viper.Set("database.password", *dbPassword)
	}
	if *dbName != "" {
		viper.Set("database.dbname", *dbName)
	}
	if *dbSSLMode != "" {
		viper.Set("database.sslmode", *dbSSLMode)
	}
	if *grinexBaseURL != "" {
		viper.Set("grinex.base_url", *grinexBaseURL)
	}
	if *grinexTimeout != "" {
		viper.Set("grinex.timeout", *grinexTimeout)
	}
	if *logLevel != "" {
		viper.Set("logging.level", *logLevel)
	}
}

func getString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
